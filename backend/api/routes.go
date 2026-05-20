package api

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/modules/wellness/application"
)

// RegisterWellnessRoutes mounts wellness HTTP routes.
func RegisterWellnessRoutes(h *server.Hertz, exec *application.Executor, guestMw app.HandlerFunc, rl *infra.RedisRateLimiter, streamPerMinute int) {
	h.POST("/v1/sessions/stream", guestMw, NewStreamHandler(exec, rl, streamPerMinute))
	h.GET("/v1/sessions/:id", guestMw, NewGetSessionHandler(exec.Store))
	h.POST("/v1/sessions/end", guestMw, NewEndSessionHandler(exec.Store))
	h.GET("/v1/profile", guestMw, NewGetProfileHandler(exec.Store))
	h.PUT("/v1/profile", guestMw, NewPutProfileHandler(exec.Store))
}
