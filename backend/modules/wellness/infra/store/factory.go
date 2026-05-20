package store

import (
	"context"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	"os"
)

// NewFromEnv returns MongoStore when USE_MEMORY_STORE is not "true", else MemoryStore.
func NewFromEnv(ctx context.Context) (domain.Store, error) {
	if os.Getenv("USE_MEMORY_STORE") == "true" {
		return NewMemoryStore(), nil
	}
	return NewMongoStore(ctx)
}
