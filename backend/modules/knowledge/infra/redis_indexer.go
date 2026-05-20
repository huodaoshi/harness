package infra

import (
	"context"

	redisindexer "github.com/cloudwego/eino-ext/components/indexer/redis"
)

// Indexer wraps eino-ext Redis indexer for knowledge pipeline.
type Indexer = redisindexer.Indexer

// IndexerConfig configures the Redis indexer.
type IndexerConfig = redisindexer.IndexerConfig

// Hashes is the Redis HASH payload for one chunk.
type Hashes = redisindexer.Hashes

// FieldValue is one field in a Redis HASH document.
type FieldValue = redisindexer.FieldValue

// NewIndexer creates a Redis-backed vector indexer.
func NewIndexer(ctx context.Context, cfg *IndexerConfig) (*Indexer, error) {
	return redisindexer.NewIndexer(ctx, cfg)
}
