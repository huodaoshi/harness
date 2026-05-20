package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/huodaoshi/harness/backend/pkg/apierror"
)

const defaultStreamRateLimitPerMinute = 60

// streamRateLimiter checks per-user stream quotas.
type streamRateLimiter interface {
	Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error)
}

// checkStreamRateLimit returns false when the client should receive 429.
func checkStreamRateLimit(ctx context.Context, c *app.RequestContext, rl streamRateLimiter, perMinute int) bool {
	if rl == nil || perMinute <= 0 {
		return true
	}

	userID, _ := UserIDFromContext(c)
	key := fmt.Sprintf("ratelimit:%s:stream", userID)
	if userID == "" {
		key = fmt.Sprintf("ratelimit:ip:%s:stream", extractIP(c))
	}

	ok, err := rl.Allow(ctx, key, int64(perMinute), time.Minute)
	if err != nil {
		slog.WarnContext(ctx, "stream: rate limit check failed", "error", err)
		return true
	}
	if ok {
		return true
	}

	apierror.Render(ctx, c, apierror.ErrRateLimit)
	return false
}
