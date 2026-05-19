package store

import (
	"context"
	"os"
)

// NewFromEnv returns MongoStore when USE_MEMORY_STORE is not "true", else MemoryStore.
func NewFromEnv(ctx context.Context) (Store, error) {
	if os.Getenv("USE_MEMORY_STORE") == "true" {
		return NewMemoryStore(), nil
	}
	return NewMongoStore(ctx)
}
