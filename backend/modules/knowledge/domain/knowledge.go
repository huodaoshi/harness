package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// ErrSpaceNotFound is returned when a requested space does not exist.
var ErrSpaceNotFound = errors.New("space not found")

const (
	IngestStatusPending       = 1
	IngestStatusProcessing    = 2
	IngestStatusDone          = 3
	IngestStatusFailed        = 4
	IngestStatusPendingReview = 5
	IngestStatusInvalidated   = 6
	IngestStatusDeleted       = 7
)

// Space represents a knowledge space (collection of documents).
type Space struct {
	SpaceID     int64      `bson:"space_id"`
	Name        string     `bson:"name"`
	Description string     `bson:"description"`
	ChunkCount  int64      `bson:"chunk_count"`
	DocTypes    []int      `bson:"doc_types"`
	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at"`
}

// Chunk represents a knowledge chunk retrieved from vector search.
type Chunk struct {
	ChunkID    string
	Content    string
	SourceType int
	DocType    int
	SourceURL  string
	Score      float64
}

// SpaceRepo defines the persistence interface for knowledge spaces.
type SpaceRepo interface {
	ListActive(ctx context.Context) ([]*Space, error)
	GetByID(ctx context.Context, spaceID int64) (*Space, error)
	AdjustChunkCount(ctx context.Context, spaceID int64, delta int64) error
}

// VectorRepo defines the vector search interface.
type VectorRepo interface {
	Search(ctx context.Context, spaceID int64, embedding []float32, topK int) ([]*Chunk, error)
}

// UpdateSpaceInput carries the partial-update fields for a space.
type UpdateSpaceInput struct {
	Name        *string
	Description *string
}

// SpaceAdminRepo defines the write-side persistence interface for knowledge spaces.
type SpaceAdminRepo interface {
	Create(ctx context.Context, s *Space) error
	Update(ctx context.Context, spaceID int64, input UpdateSpaceInput) error
	SoftDelete(ctx context.Context, spaceID int64) error
}

// SpaceAdminQueryRepo defines read-side admin queries on knowledge spaces.
type SpaceAdminQueryRepo interface {
	ListActive(ctx context.Context) ([]*Space, error)
}

// SpaceCleanupRepo handles post-deletion async cleanup for a space.
type SpaceCleanupRepo interface {
	DropVectorIndex(ctx context.Context, spaceID int64) error
	UnlinkSessions(ctx context.Context, spaceID int64) error
}

// IngestJob represents a document ingestion job.
type IngestJob struct {
	JobID      string    `bson:"job_id"       json:"job_id"`
	SpaceID    int64     `bson:"space_id"     json:"space_id"`
	SourceType int       `bson:"source_type"  json:"source_type"`
	SourceURL  string    `bson:"source_url"   json:"source_url,omitempty"`
	Content    string    `bson:"content"      json:"content,omitempty"`
	DocType    int       `bson:"doc_type"     json:"doc_type"`
	DocKey     string    `bson:"doc_key"      json:"doc_key,omitempty"`
	ChunkCount int       `bson:"chunk_count"  json:"chunk_count"`
	Status     int       `bson:"status"       json:"status"`
	Error      string    `bson:"error"        json:"error,omitempty"`
	CreatedAt  time.Time `bson:"created_at"   json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"   json:"updated_at"`
}

// IngestJobRepo defines the persistence interface for ingestion jobs.
type IngestJobRepo interface {
	Create(ctx context.Context, job *IngestJob) error
	GetByID(ctx context.Context, jobID string) (*IngestJob, error)
	UpdateStatus(ctx context.Context, jobID string, status int, errMsg string) error
	UpdateChunkCount(ctx context.Context, jobID string, count int) error
	Delete(ctx context.Context, jobID string) error
	ListBySpace(ctx context.Context, spaceID int64, status *int, page, pageSize int) ([]*IngestJob, int64, error)
	FindByDocKey(ctx context.Context, docKey string) (*IngestJob, error)
}

// MQPublisher defines the interface for publishing messages to a message queue.
type MQPublisher interface {
	Publish(ctx context.Context, topic string, msg []byte) error
}

// IngestChunkRepo abstracts Redis-level chunk lifecycle operations.
type IngestChunkRepo interface {
	MarkChunksInactive(ctx context.Context, spaceID int64, docKey string) (int, error)
	DeleteDocSetAndChunks(ctx context.Context, spaceID int64, docKey string) error
	ChunkCount(ctx context.Context, spaceID int64, docKey string) (int64, error)
}

func NormalizeIngestContent(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	return strings.TrimSpace(content)
}

func BuildIngestDocKey(spaceID int64, sourceType int, sourceURL, content string) (string, error) {
	switch sourceType {
	case 1:
		normalized := NormalizeIngestContent(content)
		sum := sha256.Sum256([]byte(normalized))
		return fmt.Sprintf("manual:%d:%s", spaceID, hex.EncodeToString(sum[:])), nil
	case 2:
		normalized, err := normalizeSourceURLForDocKey(sourceURL)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("url:%d:%s", spaceID, normalized), nil
	case 3:
		return fmt.Sprintf("file:%d:%s", spaceID, content), nil
	default:
		return "", fmt.Errorf("unsupported source_type %d", sourceType)
	}
}

func normalizeSourceURLForDocKey(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("parse source_url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid source_url")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if port != "" {
		defaultPort := (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443")
		if !defaultPort {
			host = net.JoinHostPort(host, port)
		}
	}
	u.Host = host
	u.RawQuery = ""
	u.Fragment = ""
	if u.Path == "" {
		u.Path = "/"
	}
	return u.String(), nil
}
