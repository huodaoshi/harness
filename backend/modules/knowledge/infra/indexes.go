package infra

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureKnowledgeIndexes creates MongoDB indexes for knowledge collections.
// Safe to call repeatedly (idempotent).
func EnsureKnowledgeIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection("ingest_jobs")
	model := mongo.IndexModel{
		Keys: bson.D{
			{Key: "space_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().SetName("idx_space_created"),
	}
	if _, err := coll.Indexes().CreateOne(ctx, model); err != nil && !isIndexExistsErr(err) {
		return fmt.Errorf("knowledge indexes: ingest_jobs idx_space_created: %w", err)
	}
	return nil
}

func isIndexExistsErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "index already exists") ||
		strings.Contains(s, "indexoptionsconflict") ||
		strings.Contains(s, "indexkeyspecsconflict")
}
