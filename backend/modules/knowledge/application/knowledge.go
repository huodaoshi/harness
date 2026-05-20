package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	objinfra "github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
)

// KnowledgeService defines knowledge retrieval operations.
type KnowledgeService interface {
	RetrieveChunks(ctx context.Context, spaceID *int64, query string, topK int, embedding []float32) ([]*domain.Chunk, error)
	ListSpaces(ctx context.Context) ([]*domain.Space, error)
}

// knowledgeService is the private implementation of KnowledgeService.
type knowledgeService struct {
	spaceRepo  domain.SpaceRepo
	vectorRepo domain.VectorRepo
}

// NewKnowledgeService creates a KnowledgeService.
func NewKnowledgeService(spaceRepo domain.SpaceRepo, vectorRepo domain.VectorRepo) KnowledgeService {
	return &knowledgeService{
		spaceRepo:  spaceRepo,
		vectorRepo: vectorRepo,
	}
}

// RetrieveChunks retrieves knowledge chunks for a given query embedding.
// Returns an empty list if spaceID is nil.
func (s *knowledgeService) RetrieveChunks(ctx context.Context, spaceID *int64, query string, topK int, embedding []float32) ([]*domain.Chunk, error) {
	if spaceID == nil {
		return nil, nil
	}
	chunks, err := s.vectorRepo.Search(ctx, *spaceID, embedding, topK)
	if err != nil {
		return nil, fmt.Errorf("application: knowledge: retrieve chunks: %w", err)
	}
	return chunks, nil
}

// ListSpaces returns all active knowledge spaces.
func (s *knowledgeService) ListSpaces(ctx context.Context) ([]*domain.Space, error) {
	spaces, err := s.spaceRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("application: knowledge: list spaces: %w", err)
	}
	return spaces, nil
}

// ---------------------------------------------------------------------------
// SpaceAdminService
// ---------------------------------------------------------------------------

// SpaceAdminService defines space management operations for admin users.
type SpaceAdminService interface {
	CreateSpace(ctx context.Context, name, description string, docTypes []int) (*domain.Space, error)
	UpdateSpace(ctx context.Context, spaceID int64, input domain.UpdateSpaceInput) error
	DeleteSpace(ctx context.Context, spaceID int64) error
	ListSpaces(ctx context.Context) ([]*domain.Space, error)
}

// spaceAdminService is the private implementation of SpaceAdminService.
type spaceAdminService struct {
	adminRepo   domain.SpaceAdminRepo
	cleanupRepo domain.SpaceCleanupRepo
	queryRepo   domain.SpaceAdminQueryRepo
}

// NewSpaceAdminService creates a SpaceAdminService.
func NewSpaceAdminService(adminRepo domain.SpaceAdminRepo, cleanupRepo domain.SpaceCleanupRepo, queryRepo domain.SpaceAdminQueryRepo) SpaceAdminService {
	return &spaceAdminService{
		adminRepo:   adminRepo,
		cleanupRepo: cleanupRepo,
		queryRepo:   queryRepo,
	}
}

// CreateSpace creates a new knowledge space.
func (s *spaceAdminService) CreateSpace(ctx context.Context, name, description string, docTypes []int) (*domain.Space, error) {
	space := &domain.Space{
		Name:        name,
		Description: description,
		DocTypes:    docTypes,
	}
	if err := s.adminRepo.Create(ctx, space); err != nil {
		return nil, fmt.Errorf("application: space admin: create space: %w", err)
	}
	return space, nil
}

// UpdateSpace updates the name and/or description of an existing space.
func (s *spaceAdminService) UpdateSpace(ctx context.Context, spaceID int64, input domain.UpdateSpaceInput) error {
	if err := s.adminRepo.Update(ctx, spaceID, input); err != nil {
		if errors.Is(err, domain.ErrSpaceNotFound) {
			return apierror.ErrNotFound
		}
		return fmt.Errorf("application: space admin: update space: %w", err)
	}
	return nil
}

// DeleteSpace soft-deletes a space and asynchronously cleans up its vector
// index and associated session references.
func (s *spaceAdminService) DeleteSpace(ctx context.Context, spaceID int64) error {
	if err := s.adminRepo.SoftDelete(ctx, spaceID); err != nil {
		if errors.Is(err, domain.ErrSpaceNotFound) {
			return apierror.ErrNotFound
		}
		return fmt.Errorf("application: space admin: delete space: %w", err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.cleanupRepo.DropVectorIndex(bgCtx, spaceID); err != nil {
			slog.WarnContext(bgCtx, "application: space admin: delete space", "error", err)
		}

		if err := s.cleanupRepo.UnlinkSessions(bgCtx, spaceID); err != nil {
			slog.WarnContext(bgCtx, "application: space admin: delete space", "error", err)
		}
	}()

	return nil
}

// ListSpaces returns all active (non-deleted) knowledge spaces.
func (s *spaceAdminService) ListSpaces(ctx context.Context) ([]*domain.Space, error) {
	spaces, err := s.queryRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("application: space admin: list spaces: %w", err)
	}
	return spaces, nil
}

// ---------------------------------------------------------------------------
// IngestService
// ---------------------------------------------------------------------------

// IngestService defines document ingestion operations.
type IngestService interface {
	SubmitIngestJob(ctx context.Context, spaceID int64, sourceType int, sourceURL, content string, docType int, docKey string) (*domain.IngestJob, error)
	GetIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error)
	ListIngestJobs(ctx context.Context, spaceID int64, status *int, page, pageSize int) ([]*domain.IngestJob, int64, error)
	FindJobByDocKey(ctx context.Context, docKey string) (*domain.IngestJob, error)
	RetryIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error)
	ApproveIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error)
	DeleteDraftJob(ctx context.Context, jobID string) (*domain.IngestJob, error)
	InvalidateIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error)
	HardDeleteIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error)
}

// ingestService is the private implementation of IngestService.
type ingestService struct {
	spaceRepo  domain.SpaceRepo
	jobRepo    domain.IngestJobRepo
	mq         domain.MQPublisher
	outbox     domain.MQOutboxRepo
	chunkRepo  domain.IngestChunkRepo
	objStorage objinfra.ObjectStorage
}

// NewIngestService creates an IngestService.
func NewIngestService(
	spaceRepo domain.SpaceRepo,
	jobRepo domain.IngestJobRepo,
	mq domain.MQPublisher,
	outbox domain.MQOutboxRepo,
	chunkRepo domain.IngestChunkRepo,
	objStorage objinfra.ObjectStorage,
) IngestService {
	return &ingestService{
		spaceRepo:  spaceRepo,
		jobRepo:    jobRepo,
		mq:         mq,
		outbox:     outbox,
		chunkRepo:  chunkRepo,
		objStorage: objStorage,
	}
}

const maxIngestJobErrStoredLen = 2000

func truncateIngestJobErrMsg(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxIngestJobErrStoredLen {
		return s
	}
	return s[:maxIngestJobErrStoredLen] + "..."
}

func containsDocType(docTypes []int, docType int) bool {
	for _, v := range docTypes {
		if v == docType {
			return true
		}
	}
	return false
}

func (s *ingestService) publish(ctx context.Context, job *domain.IngestJob) error {
	msgBytes, err := json.Marshal(job)
	if err != nil {
		if upErr := s.jobRepo.UpdateStatus(ctx, job.JobID, domain.IngestStatusFailed, truncateIngestJobErrMsg(err.Error())); upErr != nil {
			slog.WarnContext(ctx, "application: ingest: publish: update status failed after marshal failure",
				"job_id", job.JobID, "error", upErr)
		}
		return fmt.Errorf("application: ingest: publish: marshal job: %w", err)
	}
	if err := s.mq.Publish(ctx, "knowledge-ingest", msgBytes); err != nil {
		if upErr := s.outbox.UpsertPendingAfterPublishFailure(ctx, job.JobID, "knowledge-ingest", msgBytes); upErr != nil {
			slog.WarnContext(ctx, "application: ingest: publish: outbox upsert after publish failure",
				"job_id", job.JobID, "error", upErr)
			return fmt.Errorf("application: ingest: publish: publish job: %w; outbox: %v", err, upErr)
		}
		slog.WarnContext(ctx, "application: ingest: publish: publish failed (queued in outbox)",
			"job_id", job.JobID, "error", err)
	}
	return nil
}

func bestEffortDeleteObject(ctx context.Context, storage objinfra.ObjectStorage, key string) {
	if storage == nil || strings.TrimSpace(key) == "" {
		return
	}
	if err := storage.Delete(ctx, key); err != nil {
		slog.WarnContext(ctx, "application: ingest: best effort delete object", "key", key, "error", err)
	}
}

func bestEffortAdjustChunkCount(ctx context.Context, repo domain.SpaceRepo, spaceID int64, delta int64) {
	if repo == nil || delta == 0 {
		return
	}
	if err := repo.AdjustChunkCount(ctx, spaceID, delta); err != nil {
		slog.WarnContext(ctx, "application: ingest: best effort adjust chunk count",
			"space_id", spaceID, "delta", delta, "error", err)
	}
}

func (s *ingestService) loadJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: load job: %w", err)
	}
	if job == nil {
		return nil, apierror.ErrNotFound
	}
	return job, nil
}

// SubmitIngestJob creates an ingestion job.
// If docKey is non-empty it is used directly; otherwise it is derived from spaceID, sourceType, sourceURL and content.
func (s *ingestService) SubmitIngestJob(ctx context.Context, spaceID int64, sourceType int, sourceURL, content string, docType int, docKey string) (*domain.IngestJob, error) {
	space, err := s.spaceRepo.GetByID(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: submit: get space: %w", err)
	}
	if space == nil {
		return nil, apierror.ErrNotFound
	}
	if !containsDocType(space.DocTypes, docType) {
		return nil, apierror.ErrBadRequest
	}

	if docKey == "" {
		docKey, err = domain.BuildIngestDocKey(spaceID, sourceType, sourceURL, content)
		if err != nil {
			return nil, apierror.ErrBadRequest
		}
	}

	job := &domain.IngestJob{
		SpaceID:    spaceID,
		SourceType: sourceType,
		SourceURL:  strings.TrimSpace(sourceURL),
		Content:    content,
		DocType:    docType,
		DocKey:     docKey,
		Status:     domain.IngestStatusPending,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("application: ingest: submit: create job: %w", err)
	}
	if err := s.publish(ctx, job); err != nil {
		return nil, fmt.Errorf("application: ingest: submit: publish: %w", err)
	}

	return job, nil
}

// GetIngestJob retrieves an ingestion job by ID.
func (s *ingestService) GetIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: get job: %w", err)
	}
	if job == nil {
		return nil, apierror.ErrNotFound
	}
	return job, nil
}

// ListIngestJobs returns a paginated list of ingestion jobs for a given space.
func (s *ingestService) ListIngestJobs(ctx context.Context, spaceID int64, status *int, page, pageSize int) ([]*domain.IngestJob, int64, error) {
	space, err := s.spaceRepo.GetByID(ctx, spaceID)
	if err != nil {
		return nil, 0, fmt.Errorf("application: ingest: list jobs: get space: %w", err)
	}
	if space == nil {
		return nil, 0, apierror.ErrNotFound
	}
	jobs, total, err := s.jobRepo.ListBySpace(ctx, spaceID, status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("application: ingest: list jobs: %w", err)
	}
	return jobs, total, nil
}

// FindJobByDocKey returns the most recent ingestion job matching docKey, or nil if none found.
func (s *ingestService) FindJobByDocKey(ctx context.Context, docKey string) (*domain.IngestJob, error) {
	job, err := s.jobRepo.FindByDocKey(ctx, docKey)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: find by doc_key: %w", err)
	}
	return job, nil
}

// RetryIngestJob transitions a Failed job back to Pending and re-publishes it to the MQ.
func (s *ingestService) RetryIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: retry: %w", err)
	}
	if job.Status != domain.IngestStatusFailed {
		return nil, apierror.ErrBadRequest
	}
	if err := s.jobRepo.UpdateStatus(ctx, job.JobID, domain.IngestStatusPending, ""); err != nil {
		return nil, fmt.Errorf("application: ingest: retry: update status: %w", err)
	}
	job.Status = domain.IngestStatusPending
	if err := s.publish(ctx, job); err != nil {
		return nil, fmt.Errorf("application: ingest: retry: %w", err)
	}
	return job, nil
}

// ApproveIngestJob transitions a PendingReview job to Pending and publishes it to the MQ.
func (s *ingestService) ApproveIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: approve: %w", err)
	}
	if job.Status != domain.IngestStatusPendingReview {
		return nil, apierror.ErrBadRequest
	}
	if err := s.jobRepo.UpdateStatus(ctx, job.JobID, domain.IngestStatusPending, ""); err != nil {
		return nil, fmt.Errorf("application: ingest: approve: update status: %w", err)
	}
	job.Status = domain.IngestStatusPending
	if err := s.publish(ctx, job); err != nil {
		return nil, fmt.Errorf("application: ingest: approve: %w", err)
	}
	return job, nil
}

// DeleteDraftJob removes a job that is still in Failed or PendingReview state (never indexed).
// For file uploads (source_type=3) the object is deleted best-effort.
func (s *ingestService) DeleteDraftJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: delete draft: %w", err)
	}
	if job.Status != domain.IngestStatusFailed && job.Status != domain.IngestStatusPendingReview {
		return nil, apierror.ErrBadRequest
	}
	if job.SourceType == 3 {
		bestEffortDeleteObject(ctx, s.objStorage, job.SourceURL)
	}
	if err := s.jobRepo.Delete(ctx, job.JobID); err != nil {
		return nil, fmt.Errorf("application: ingest: delete draft: %w", err)
	}
	return job, nil
}

// InvalidateIngestJob marks all indexed chunks of a Done job as inactive,
// deducts their count from Space.chunk_count, and transitions the job to Invalidated.
func (s *ingestService) InvalidateIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: invalidate: %w", err)
	}
	if job.Status != domain.IngestStatusDone {
		return nil, apierror.ErrBadRequest
	}

	markedCount, markErr := s.chunkRepo.MarkChunksInactive(ctx, job.SpaceID, job.DocKey)
	if markErr != nil {
		return nil, fmt.Errorf("application: ingest: invalidate: mark inactive: %w", markErr)
	}

	chunkDelta := int64(job.ChunkCount)
	if chunkDelta == 0 {
		cnt, scardErr := s.chunkRepo.ChunkCount(ctx, job.SpaceID, job.DocKey)
		if scardErr != nil {
			slog.WarnContext(ctx, "application: ingest: invalidate: best effort chunk count", "space_id", job.SpaceID, "error", scardErr)
			chunkDelta = int64(markedCount)
		} else {
			chunkDelta = cnt
		}
	}
	bestEffortAdjustChunkCount(ctx, s.spaceRepo, job.SpaceID, -chunkDelta)

	if err := s.jobRepo.UpdateStatus(ctx, job.JobID, domain.IngestStatusInvalidated, ""); err != nil {
		return nil, fmt.Errorf("application: ingest: invalidate: update status: %w", err)
	}
	job.Status = domain.IngestStatusInvalidated
	return job, nil
}

// HardDeleteIngestJob deletes all Redis chunks of an Invalidated job, best-effort deletes the file,
// and transitions the job to Deleted.
func (s *ingestService) HardDeleteIngestJob(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	job, err := s.loadJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("application: ingest: hard delete: %w", err)
	}
	if job.Status != domain.IngestStatusInvalidated {
		return nil, apierror.ErrBadRequest
	}

	if err := s.chunkRepo.DeleteDocSetAndChunks(ctx, job.SpaceID, job.DocKey); err != nil {
		return nil, fmt.Errorf("application: ingest: hard delete: delete chunks: %w", err)
	}

	if job.SourceType == 3 {
		bestEffortDeleteObject(ctx, s.objStorage, job.SourceURL)
	}

	if err := s.jobRepo.UpdateStatus(ctx, job.JobID, domain.IngestStatusDeleted, ""); err != nil {
		return nil, fmt.Errorf("application: ingest: hard delete: update status: %w", err)
	}
	job.Status = domain.IngestStatusDeleted
	return job, nil
}
