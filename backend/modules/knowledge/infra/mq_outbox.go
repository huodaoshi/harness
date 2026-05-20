package infra

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

// MongoMQOutboxRepo implements domain.MQOutboxRepo.
type MongoMQOutboxRepo struct {
	coll    *mongo.Collection
	idg     *idgen.IDGenerator
	idxOnce sync.Once
	idxErr  error
}

// NewMongoMQOutboxRepo creates a MongoMQOutboxRepo.
func NewMongoMQOutboxRepo(db *mongo.Database, idg *idgen.IDGenerator) *MongoMQOutboxRepo {
	return &MongoMQOutboxRepo{
		coll: db.Collection("mq_outbox"),
		idg:  idg,
	}
}

func (r *MongoMQOutboxRepo) ensureIndexes(ctx context.Context) error {
	r.idxOnce.Do(func() {
		_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "job_id", Value: 1}, {Key: "topic", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_job_topic"),
		})
		if err != nil && !isMongoIndexAlreadyExists(err) {
			r.idxErr = err
		}
	})
	return r.idxErr
}

func isMongoIndexAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "already exists") ||
		strings.Contains(s, "IndexOptionsConflict") ||
		strings.Contains(s, "duplicate key")
}

// UpsertPendingAfterPublishFailure upserts by (job_id, topic).
func (r *MongoMQOutboxRepo) UpsertPendingAfterPublishFailure(ctx context.Context, jobID, topic string, payload []byte) error {
	if err := r.ensureIndexes(ctx); err != nil {
		return fmt.Errorf("infra: mq outbox: indexes: %w", err)
	}
	now := time.Now()
	filter := bson.M{"job_id": jobID, "topic": topic}
	outboxID := r.idg.Generate()
	update := bson.M{
		"$set": bson.M{
			"payload":    payload,
			"updated_at": now,
			"last_error": "",
		},
		"$setOnInsert": bson.M{
			"outbox_id":  outboxID,
			"attempts":   0,
			"created_at": now,
		},
	}
	opts := options.Update().SetUpsert(true)
	if _, err := r.coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("infra: mq outbox: upsert: %w", err)
	}
	return nil
}

// ListPending returns the oldest pending outbox rows (all rows are pending until deleted).
func (r *MongoMQOutboxRepo) ListPending(ctx context.Context, limit int64) ([]*domain.MQOutbox, error) {
	if limit <= 0 {
		limit = 32
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}).SetLimit(limit)
	cur, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("infra: mq outbox: list: %w", err)
	}
	defer cur.Close(ctx)

	var out []*domain.MQOutbox
	if err = cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("infra: mq outbox: list decode: %w", err)
	}
	return out, nil
}

// Delete removes one outbox document by outbox_id.
func (r *MongoMQOutboxRepo) Delete(ctx context.Context, outboxID string) error {
	if _, err := r.coll.DeleteOne(ctx, bson.M{"outbox_id": outboxID}); err != nil {
		return fmt.Errorf("infra: mq outbox: delete: %w", err)
	}
	return nil
}

const maxOutboxErrLen = 1500

func truncateOutboxErr(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxOutboxErrLen {
		return s
	}
	return s[:maxOutboxErrLen] + "…"
}

// IncrementFailure increments attempts and returns the new attempts count.
func (r *MongoMQOutboxRepo) IncrementFailure(ctx context.Context, outboxID string, errMsg string) (int, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.M{"outbox_id": outboxID}
	update := bson.M{
		"$inc": bson.M{"attempts": 1},
		"$set": bson.M{
			"last_error": truncateOutboxErr(errMsg),
			"updated_at": time.Now(),
		},
	}
	var doc domain.MQOutbox
	sr := r.coll.FindOneAndUpdate(ctx, filter, update, opts)
	if err := sr.Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, fmt.Errorf("infra: mq outbox: increment failure: %w", err)
	}
	return doc.Attempts, nil
}

// Ensure implementation satisfies domain.
var _ domain.MQOutboxRepo = (*MongoMQOutboxRepo)(nil)
