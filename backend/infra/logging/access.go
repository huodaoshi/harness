package logging

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-Id"

// AccessLogMiddleware logs one line per HTTP request (method, path, status, latency).
func AccessLogMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		reqID := strings.TrimSpace(string(c.GetHeader(requestIDHeader)))
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Response.Header.Set(requestIDHeader, reqID)
		c.Set("request_id", reqID)

		c.Next(ctx)

		status := c.Response.StatusCode()
		if status == 0 {
			status = 200
		}
		attrs := []any{
			"request_id", reqID,
			"method", string(c.Method()),
			"path", string(c.Path()),
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", clientIP(c),
		}
		if query := string(c.URI().QueryString()); query != "" {
			attrs = append(attrs, "query", query)
		}
		if status >= 500 {
			slog.ErrorContext(ctx, "http request", attrs...)
		} else if status >= 400 {
			slog.WarnContext(ctx, "http request", attrs...)
		} else {
			slog.InfoContext(ctx, "http request", attrs...)
		}
	}
}

func clientIP(c *app.RequestContext) string {
	if xff := string(c.GetHeader("X-Forwarded-For")); xff != "" {
		if i := strings.Index(xff, ","); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if xri := string(c.GetHeader("X-Real-IP")); xri != "" {
		return strings.TrimSpace(xri)
	}
	return c.ClientIP()
}
