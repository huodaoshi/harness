package infra

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func hashDocKey(docKey string) string {
	sum := sha256.Sum256([]byte(docKey))
	return hex.EncodeToString(sum[:])
}

func BuildDocSetKey(spaceID int64, docKey string) string {
	return fmt.Sprintf("knowledge_docset:%d:%s", spaceID, hashDocKey(docKey))
}

func BuildDocStageKey(spaceID int64, docKey, jobID string) string {
	return fmt.Sprintf("knowledge_docset_staging:%d:%s:%s", spaceID, hashDocKey(docKey), jobID)
}

func BuildDocLockKey(spaceID int64, docKey string) string {
	return fmt.Sprintf("knowledge_doclock:%d:%s", spaceID, hashDocKey(docKey))
}

func AcquireDocLock(ctx context.Context, rdb redis.UniversalClient, spaceID int64, docKey, owner string, ttl time.Duration) (bool, error) {
	ok, err := rdb.SetNX(ctx, BuildDocLockKey(spaceID, docKey), owner, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("infra: vector docset: acquire lock: %w", err)
	}
	return ok, nil
}

func ReleaseDocLock(ctx context.Context, rdb redis.UniversalClient, spaceID int64, docKey, owner string) error {
	lockKey := BuildDocLockKey(spaceID, docKey)
	current, err := rdb.Get(ctx, lockKey).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("infra: vector docset: get lock: %w", err)
	}
	if current != owner {
		return nil
	}
	if err := rdb.Del(ctx, lockKey).Err(); err != nil {
		return fmt.Errorf("infra: vector docset: release lock: %w", err)
	}
	return nil
}

func AddStageChunkKeys(ctx context.Context, rdb redis.UniversalClient, stageKey string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	args := make([]interface{}, 0, len(keys)+1)
	args = append(args, stageKey)
	for _, key := range keys {
		args = append(args, key)
	}
	if err := rdb.SAdd(ctx, stageKey, args[1:]...).Err(); err != nil {
		return fmt.Errorf("infra: vector docset: add stage chunk keys: %w", err)
	}
	return nil
}

func LoadLiveChunkKeys(ctx context.Context, rdb redis.UniversalClient, liveKey string) ([]string, error) {
	keys, err := rdb.SMembers(ctx, liveKey).Result()
	if err != nil {
		return nil, fmt.Errorf("infra: vector docset: load live chunk keys: %w", err)
	}
	return keys, nil
}

func PromoteStageToLive(ctx context.Context, rdb redis.UniversalClient, stageKey, liveKey string) error {
	keys, err := rdb.SMembers(ctx, stageKey).Result()
	if err != nil {
		return fmt.Errorf("infra: vector docset: read stage members: %w", err)
	}
	pipe := rdb.TxPipeline()
	pipe.Del(ctx, liveKey)
	if len(keys) > 0 {
		args := make([]interface{}, 0, len(keys))
		for _, key := range keys {
			args = append(args, key)
		}
		pipe.SAdd(ctx, liveKey, args...)
	}
	pipe.Del(ctx, stageKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("infra: vector docset: promote stage to live: %w", err)
	}
	return nil
}

func DeleteChunkKeys(ctx context.Context, rdb redis.UniversalClient, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	args := make([]string, 0, len(keys))
	args = append(args, keys...)
	if err := rdb.Del(ctx, args...).Err(); err != nil {
		return fmt.Errorf("infra: vector docset: delete chunk keys: %w", err)
	}
	return nil
}

func DeleteStageArtifacts(ctx context.Context, rdb redis.UniversalClient, stageKey string, stageChunkKeys []string) error {
	pipe := rdb.TxPipeline()
	if len(stageChunkKeys) > 0 {
		args := make([]string, 0, len(stageChunkKeys))
		args = append(args, stageChunkKeys...)
		pipe.Del(ctx, args...)
	}
	pipe.Del(ctx, stageKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("infra: vector docset: delete stage artifacts: %w", err)
	}
	return nil
}

func MarkChunksInactive(ctx context.Context, rdb redis.UniversalClient, docSetKey string) (int, error) {
	keys, err := rdb.SMembers(ctx, docSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("infra: vector docset: mark inactive: load members: %w", err)
	}
	if len(keys) == 0 {
		return 0, nil
	}

	pipe := rdb.Pipeline()
	for _, key := range keys {
		pipe.HSet(ctx, key, "is_active", "invalid")
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("infra: vector docset: mark inactive: exec: %w", err)
	}
	return len(keys), nil
}

func DeleteDocSetAndChunks(ctx context.Context, rdb redis.UniversalClient, docSetKey string) error {
	keys, err := rdb.SMembers(ctx, docSetKey).Result()
	if err != nil {
		return fmt.Errorf("infra: vector docset: delete doc set: load members: %w", err)
	}

	pipe := rdb.TxPipeline()
	if len(keys) > 0 {
		args := make([]string, 0, len(keys))
		args = append(args, keys...)
		pipe.Del(ctx, args...)
	}
	pipe.Del(ctx, docSetKey)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("infra: vector docset: delete doc set: exec: %w", err)
	}
	return nil
}
