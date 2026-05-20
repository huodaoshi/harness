package infra

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

// ---------------------------------------------------------------------------
// MongoSpaceRepo
// ---------------------------------------------------------------------------

// MongoSpaceRepo implements domain.SpaceRepo backed by MongoDB.
type MongoSpaceRepo struct {
	coll *mongo.Collection
}

// NewMongoSpaceRepo creates a MongoSpaceRepo.
func NewMongoSpaceRepo(db *mongo.Database) domain.SpaceRepo {
	return &MongoSpaceRepo{coll: db.Collection("spaces")}
}

// ListActive returns all spaces that have not been soft-deleted.
func (r *MongoSpaceRepo) ListActive(ctx context.Context) ([]*domain.Space, error) {
	filter := bson.M{"deleted_at": nil}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("infra: space: list active: find: %w", err)
	}
	defer cursor.Close(ctx)

	var spaces []*domain.Space
	if err = cursor.All(ctx, &spaces); err != nil {
		return nil, fmt.Errorf("infra: space: list active: decode: %w", err)
	}
	return spaces, nil
}

// GetByID fetches a space by its numeric ID.
func (r *MongoSpaceRepo) GetByID(ctx context.Context, spaceID int64) (*domain.Space, error) {
	var s domain.Space
	if err := r.coll.FindOne(ctx, bson.M{"space_id": spaceID, "deleted_at": nil}).Decode(&s); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: space: get by id: %w", err)
	}
	return &s, nil
}

func (r *MongoSpaceRepo) AdjustChunkCount(ctx context.Context, spaceID int64, delta int64) error {
	if delta == 0 {
		return nil
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"space_id": spaceID, "deleted_at": nil}, bson.M{
		"$inc": bson.M{"chunk_count": delta},
		"$set": bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		return fmt.Errorf("infra: space: adjust chunk count: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("infra: space: adjust chunk count: %w", domain.ErrSpaceNotFound)
	}
	return nil
}

// ---------------------------------------------------------------------------
// RedisVectorRepo
// ---------------------------------------------------------------------------

// RedisVectorRepo implements domain.VectorRepo using Redis vector search.
type RedisVectorRepo struct {
	rdb redis.UniversalClient
}

// NewRedisVectorRepo creates a RedisVectorRepo.
func NewRedisVectorRepo(rdb redis.UniversalClient) domain.VectorRepo {
	return &RedisVectorRepo{rdb: rdb}
}

// Search performs a KNN vector search against the space index in Redis.
func (r *RedisVectorRepo) Search(ctx context.Context, spaceID int64, embedding []float32, topK int) ([]*domain.Chunk, error) {
	indexName := fmt.Sprintf("space_%d", spaceID)
	blob := float32SliceToBytes(embedding)
	query := fmt.Sprintf("(-@is_active:{invalid})=>[KNN %d @embedding $vec AS score]", topK)

	result, err := r.rdb.Do(ctx, "FT.SEARCH", indexName, query,
		"PARAMS", "2", "vec", blob,
		"SORTBY", "score",
		"LIMIT", "0", strconv.Itoa(topK),
		"RETURN", "4", "content", "source_type", "doc_type", "source_url",
		"DIALECT", "2",
	).Result()
	if err != nil {
		return nil, fmt.Errorf("infra: vector: search: %w", err)
	}

	return parseRedisSearchResult(result)
}

// float32SliceToBytes converts a []float32 to a little-endian raw byte slice
// suitable for Redis vector BLOB parameters.
func float32SliceToBytes(vals []float32) []byte {
	buf := make([]byte, len(vals)*4)
	for i, v := range vals {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(v))
	}
	return buf
}

// parseRedisSearchResult parses a raw FT.SEARCH result into domain chunks.
func parseRedisSearchResult(raw interface{}) ([]*domain.Chunk, error) {
	items, ok := raw.([]interface{})
	if !ok || len(items) < 1 {
		return nil, nil
	}

	var chunks []*domain.Chunk
	for i := 1; i+1 < len(items); i += 2 {
		chunkID, _ := items[i].(string)
		fields, _ := items[i+1].([]interface{})

		chunk := &domain.Chunk{ChunkID: chunkID}
		for j := 0; j+1 < len(fields); j += 2 {
			key, _ := fields[j].(string)
			val, _ := fields[j+1].(string)
			switch key {
			case "content":
				chunk.Content = val
			case "source_type":
				st, _ := strconv.Atoi(val)
				chunk.SourceType = st
			case "doc_type":
				dt, _ := strconv.Atoi(val)
				chunk.DocType = dt
			case "source_url":
				chunk.SourceURL = val
			case "score":
				sc, _ := strconv.ParseFloat(val, 64)
				chunk.Score = sc
			}
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

// Ensure implementations satisfy interfaces.
var (
	_ domain.SpaceRepo           = (*MongoSpaceRepo)(nil)
	_ domain.VectorRepo          = (*RedisVectorRepo)(nil)
	_ domain.SpaceAdminRepo      = (*MongoSpaceAdminRepo)(nil)
	_ domain.SpaceAdminQueryRepo = (*MongoSpaceAdminRepo)(nil)
	_ domain.IngestJobRepo       = (*MongoIngestJobRepo)(nil)
	_ domain.SpaceCleanupRepo    = (*SpaceCleanupRepoImpl)(nil)
	_ domain.IngestChunkRepo     = (*RedisIngestChunkRepo)(nil)
)

// ---------------------------------------------------------------------------
// RedisIngestChunkRepo
// ---------------------------------------------------------------------------

// RedisIngestChunkRepo implements domain.IngestChunkRepo using Redis.
type RedisIngestChunkRepo struct {
	rdb redis.UniversalClient
}

// NewRedisIngestChunkRepo creates a RedisIngestChunkRepo.
func NewRedisIngestChunkRepo(rdb redis.UniversalClient) domain.IngestChunkRepo {
	return &RedisIngestChunkRepo{rdb: rdb}
}

// MarkChunksInactive marks all chunks of the given doc as inactive.
func (r *RedisIngestChunkRepo) MarkChunksInactive(ctx context.Context, spaceID int64, docKey string) (int, error) {
	key := BuildDocSetKey(spaceID, docKey)
	count, err := MarkChunksInactive(ctx, r.rdb, key)
	if err != nil {
		return 0, fmt.Errorf("infra: ingest chunk: mark inactive: %w", err)
	}
	return count, nil
}

// DeleteDocSetAndChunks removes all chunk keys and the doc-set key.
func (r *RedisIngestChunkRepo) DeleteDocSetAndChunks(ctx context.Context, spaceID int64, docKey string) error {
	key := BuildDocSetKey(spaceID, docKey)
	if err := DeleteDocSetAndChunks(ctx, r.rdb, key); err != nil {
		return fmt.Errorf("infra: ingest chunk: delete doc set: %w", err)
	}
	return nil
}

// ChunkCount returns the number of live chunk keys tracked for the given doc.
func (r *RedisIngestChunkRepo) ChunkCount(ctx context.Context, spaceID int64, docKey string) (int64, error) {
	key := BuildDocSetKey(spaceID, docKey)
	cnt, err := r.rdb.SCard(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("infra: ingest chunk: chunk count: %w", err)
	}
	return cnt, nil
}

// ---------------------------------------------------------------------------
// MongoSpaceAdminRepo
// ---------------------------------------------------------------------------

// MongoSpaceAdminRepo implements domain.SpaceAdminRepo backed by MongoDB.
type MongoSpaceAdminRepo struct {
	coll    *mongo.Collection
	counter *idgen.MongoCounter
	idg     *idgen.IDGenerator
}

// NewMongoSpaceAdminRepo creates a MongoSpaceAdminRepo.
func NewMongoSpaceAdminRepo(db *mongo.Database, counter *idgen.MongoCounter, idg *idgen.IDGenerator) *MongoSpaceAdminRepo {
	return &MongoSpaceAdminRepo{
		coll:    db.Collection("spaces"),
		counter: counter,
		idg:     idg,
	}
}

// Create inserts a new Space, generating a monotonic spaceID.
func (r *MongoSpaceAdminRepo) Create(ctx context.Context, s *domain.Space) error {
	spaceID, err := r.counter.Next(ctx, idgen.CounterSpaceID)
	if err != nil {
		return fmt.Errorf("infra: space admin: create: get space id: %w", err)
	}
	s.SpaceID = spaceID
	s.CreatedAt = time.Now()

	if _, err = r.coll.InsertOne(ctx, s); err != nil {
		return fmt.Errorf("infra: space admin: create: insert: %w", err)
	}
	return nil
}

// Update sets name and/or description on a non-deleted space (partial update).
func (r *MongoSpaceAdminRepo) Update(ctx context.Context, spaceID int64, input domain.UpdateSpaceInput) error {
	setFields := bson.M{"updated_at": time.Now()}
	if input.Name != nil {
		setFields["name"] = *input.Name
	}
	if input.Description != nil {
		setFields["description"] = *input.Description
	}
	filter := bson.M{"space_id": spaceID, "deleted_at": nil}
	update := bson.M{"$set": setFields}
	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("infra: space admin: update: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("infra: space admin: update: %w", domain.ErrSpaceNotFound)
	}
	return nil
}

// SoftDelete marks a space as deleted.
func (r *MongoSpaceAdminRepo) SoftDelete(ctx context.Context, spaceID int64) error {
	now := time.Now()
	filter := bson.M{"space_id": spaceID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"deleted_at": now,
		"updated_at": now,
	}}
	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("infra: space admin: soft delete: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("infra: space admin: soft delete: %w", domain.ErrSpaceNotFound)
	}
	return nil
}

// ListActive returns all non-deleted spaces, ordered by created_at ascending.
func (r *MongoSpaceAdminRepo) ListActive(ctx context.Context) ([]*domain.Space, error) {
	filter := bson.M{"deleted_at": nil}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("infra: space admin: list active: find: %w", err)
	}
	defer cursor.Close(ctx)
	var spaces []*domain.Space
	if err = cursor.All(ctx, &spaces); err != nil {
		return nil, fmt.Errorf("infra: space admin: list active: decode: %w", err)
	}
	return spaces, nil
}

// ---------------------------------------------------------------------------
// MongoIngestJobRepo
// ---------------------------------------------------------------------------

// MongoIngestJobRepo implements domain.IngestJobRepo backed by MongoDB.
type MongoIngestJobRepo struct {
	coll *mongo.Collection
	idg  *idgen.IDGenerator
}

// NewMongoIngestJobRepo creates a MongoIngestJobRepo.
func NewMongoIngestJobRepo(db *mongo.Database, idg *idgen.IDGenerator) *MongoIngestJobRepo {
	return &MongoIngestJobRepo{
		coll: db.Collection("ingest_jobs"),
		idg:  idg,
	}
}

// Create inserts a new IngestJob, generating a JobID.
func (r *MongoIngestJobRepo) Create(ctx context.Context, job *domain.IngestJob) error {
	job.JobID = r.idg.Generate()
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now

	if _, err := r.coll.InsertOne(ctx, job); err != nil {
		return fmt.Errorf("infra: ingest job: create: %w", err)
	}
	return nil
}

// GetByID retrieves an IngestJob by JobID. Returns nil, nil if not found.
func (r *MongoIngestJobRepo) GetByID(ctx context.Context, jobID string) (*domain.IngestJob, error) {
	var job domain.IngestJob
	if err := r.coll.FindOne(ctx, bson.M{"job_id": jobID}).Decode(&job); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: ingest job: get by id: %w", err)
	}
	return &job, nil
}

// UpdateStatus sets ingest_jobs.status and error message for observability (Worker / Admin).
func (r *MongoIngestJobRepo) UpdateStatus(ctx context.Context, jobID string, status int, errMsg string) error {
	filter := bson.M{"job_id": jobID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"error":      errMsg,
			"updated_at": time.Now(),
		},
	}
	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("infra: ingest job: update status: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("infra: ingest job: update status: job not found: %s", jobID)
	}
	return nil
}

// UpdateChunkCount sets ingest_jobs.chunk_count for the given job.
func (r *MongoIngestJobRepo) UpdateChunkCount(ctx context.Context, jobID string, count int) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"job_id": jobID}, bson.M{
		"$set": bson.M{
			"chunk_count": count,
			"updated_at":  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("infra: ingest job: update chunk count: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("infra: ingest job: update chunk count: job not found: %s", jobID)
	}
	return nil
}

// Delete removes an ingest job by JobID.
func (r *MongoIngestJobRepo) Delete(ctx context.Context, jobID string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"job_id": jobID})
	if err != nil {
		return fmt.Errorf("infra: ingest job: delete: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("infra: ingest job: delete: job not found: %s", jobID)
	}
	return nil
}

// ListBySpace returns a paginated list of IngestJobs for the given space.
// status nil = all statuses. page is 1-based.
func (r *MongoIngestJobRepo) ListBySpace(ctx context.Context, spaceID int64, status *int, page, pageSize int) ([]*domain.IngestJob, int64, error) {
	filter := bson.M{"space_id": spaceID}
	if status != nil {
		filter["status"] = *status
	}
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("infra: ingest job: list by space: count: %w", err)
	}
	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(pageSize))
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("infra: ingest job: list by space: find: %w", err)
	}
	defer cursor.Close(ctx)
	var jobs []*domain.IngestJob
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, 0, fmt.Errorf("infra: ingest job: list by space: decode: %w", err)
	}
	if jobs == nil {
		jobs = []*domain.IngestJob{}
	}
	return jobs, total, nil
}

// FindByDocKey returns the most recent IngestJob matching docKey (sort: created_at:-1, limit:1).
// Returns nil, nil when not found.
func (r *MongoIngestJobRepo) FindByDocKey(ctx context.Context, docKey string) (*domain.IngestJob, error) {
	opts := options.FindOne().
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	var job domain.IngestJob
	if err := r.coll.FindOne(ctx, bson.M{"doc_key": docKey}, opts).Decode(&job); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: ingest job: find by doc_key: %w", err)
	}
	return &job, nil
}

// EnsureIngestJobDocKeyIndex creates a sparse index on doc_key in ingest_jobs.
// Safe to call repeatedly (idempotent).
func EnsureIngestJobDocKeyIndex(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection("ingest_jobs")
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "doc_key", Value: 1}},
		Options: options.Index().SetName("idx_doc_key").SetSparse(true),
	}
	_, err := coll.Indexes().CreateOne(ctx, idx)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("infra: ingest job: ensure doc_key index: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// SpaceCleanupRepoImpl
// ---------------------------------------------------------------------------

// SpaceCleanupRepoImpl implements domain.SpaceCleanupRepo.
type SpaceCleanupRepoImpl struct {
	rdb redis.UniversalClient
	db  *mongo.Database
}

// NewSpaceCleanupRepo creates a SpaceCleanupRepoImpl.
func NewSpaceCleanupRepo(rdb redis.UniversalClient, db *mongo.Database) domain.SpaceCleanupRepo {
	return &SpaceCleanupRepoImpl{rdb: rdb, db: db}
}

// DropVectorIndex deletes the Redis vector search index for the given space.
func (r *SpaceCleanupRepoImpl) DropVectorIndex(ctx context.Context, spaceID int64) error {
	indexName := fmt.Sprintf("space_%d", spaceID)
	if err := r.rdb.Do(ctx, "FT.DROPINDEX", indexName, "DD").Err(); err != nil {
		return fmt.Errorf("infra: space cleanup: drop vector index %q: %w", indexName, err)
	}
	return nil
}

// UnlinkSessions sets space_id to null for all sessions belonging to the space.
func (r *SpaceCleanupRepoImpl) UnlinkSessions(ctx context.Context, spaceID int64) error {
	filter := bson.M{"space_id": spaceID}
	update := bson.M{"$set": bson.M{"space_id": nil}}
	if _, err := r.db.Collection("sessions").UpdateMany(ctx, filter, update); err != nil {
		return fmt.Errorf("infra: space cleanup: unlink sessions for space %d: %w", spaceID, err)
	}
	return nil
}
