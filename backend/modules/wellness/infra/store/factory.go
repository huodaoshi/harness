package store

import (
	"context"
	"os"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
)

// NewFromConfig returns MemoryStore when wellness.use_memory_store or USE_MEMORY_STORE is set.
func NewFromConfig(ctx context.Context, app *conf.Config) (domain.Store, error) {
	if app != nil && app.Wellness.UseMemoryStore {
		return NewMemoryStore(), nil
	}
	if os.Getenv("USE_MEMORY_STORE") == "true" {
		return NewMemoryStore(), nil
	}
	return NewMongoStore(ctx)
}

// NewFromEnv loads conf and picks store implementation (tests).
func NewFromEnv(ctx context.Context) (domain.Store, error) {
	c, err := conf.Load()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(ctx, c)
}
