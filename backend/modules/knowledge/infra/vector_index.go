package infra

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

// EnsureKnowledgeVectorIndex creates the RedisSearch index space_{spaceID} for
// HASH documents with prefix space_{spaceID}: if it does not already exist.
// Schema matches RedisVectorRepo.Search (KNN on @embedding, RETURN fields).
func EnsureKnowledgeVectorIndex(ctx context.Context, rdb redis.UniversalClient, spaceID int64, vectorDim int) error {
	indexName := fmt.Sprintf("space_%d", spaceID)
	_, infoErr := rdb.Do(ctx, "FT.INFO", indexName).Result()
	if infoErr == nil {
		return nil
	}
	if !isUnknownIndexErr(infoErr) {
		return fmt.Errorf("infra: vector index: info %q: %w", indexName, infoErr)
	}

	prefix := indexName + ":"
	args := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", prefix,
		"SCHEMA",
		"content", "TEXT",
		"embedding", "VECTOR", "HNSW", "6", "TYPE", "FLOAT32", "DIM", vectorDim, "DISTANCE_METRIC", "COSINE",
		"source_type", "NUMERIC",
		"doc_type", "NUMERIC",
		"source_url", "TEXT",
		"is_active", "TAG",
	}
	if err := rdb.Do(ctx, args...).Err(); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "index already exists") {
			return nil
		}
		return fmt.Errorf("infra: vector index: create %q dim=%d: %w", indexName, vectorDim, err)
	}
	return nil
}

func isUnknownIndexErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "unknown index") || strings.Contains(s, "no such index")
}

func BuildKnowledgeChunkKey(spaceID int64, jobID string, chunkIndex int) string {
	return fmt.Sprintf("space_%d:chunk:%s:%06d", spaceID, jobID, chunkIndex)
}

// UpsertKnowledgeChunk writes one HASH document to the provided Redis key.
func UpsertKnowledgeChunk(
	ctx context.Context,
	rdb redis.UniversalClient,
	key string,
	content string,
	embedding []float32,
	sourceType, docType int,
	sourceURL, docKey, jobID string,
	chunkIndex int,
) error {
	blob := float32SliceToBytes(embedding)
	fields := []interface{}{
		"HSET", key,
		"content", content,
		"embedding", blob,
		"source_type", strconv.Itoa(sourceType),
		"doc_type", strconv.Itoa(docType),
		"source_url", sourceURL,
		"doc_key", docKey,
		"job_id", jobID,
		"chunk_index", strconv.Itoa(chunkIndex),
		"is_active", "active",
	}
	if err := rdb.Do(ctx, fields...).Err(); err != nil {
		return fmt.Errorf("infra: vector index: hset %q: %w", key, err)
	}
	return nil
}

// Float64SliceToFloat32 converts embedding model output to []float32 for Redis blob encoding.
func Float64SliceToFloat32(v []float64) []float32 {
	out := make([]float32, len(v))
	for i, x := range v {
		out[i] = float32(x)
	}
	return out
}
