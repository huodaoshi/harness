package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/api"
	"github.com/huodaoshi/harness/backend/api/nextchat"
	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/infra/logging"
	"github.com/huodaoshi/harness/backend/knowledgebootstrap"
	"github.com/huodaoshi/harness/backend/modules/wellness/application"
)

func main() {
	ctx := context.Background()

	bundle, err := loadConfigAndInfra(ctx)
	if err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded", "env", bundle.Cfg.App.Env, "port", bundle.Cfg.App.Port)

	auth, err := wireAuth(bundle.Cfg, bundle.MongoClient, bundle.RedisClient)
	if err != nil {
		slog.Error("auth wiring failed", "error", err)
		os.Exit(1)
	}

	exec, err := application.NewExecutor(ctx, bundle.Cfg)
	if err != nil {
		slog.Error("session executor failed", "error", err)
		os.Exit(1)
	}

	addr := listenAddr(bundle.Cfg)

	h := server.Default(server.WithHostPorts(addr))
	h.Use(logging.AccessLogMiddleware())
	api.RegisterAuthRoutes(h, api.NewAuthHandler(auth.Service), api.JWTAuthMiddleware(auth.Signer))
	streamRL := infra.NewRedisRateLimiter(bundle.RedisClient)
	api.RegisterWellnessRoutes(h, exec, api.JWTOrGuestMiddleware(auth.Signer), streamRL, bundle.Cfg.RateLimit.StreamPerMinute)
	nextchat.Register(h, bundle.Cfg)

	kb, err := knowledgebootstrap.Wire(ctx, bundle.Cfg, bundle.MongoClient, bundle.RedisClient, true)
	if err != nil {
		slog.Error("knowledge wiring failed", "error", err)
		os.Exit(1)
	}
	defer kb.Close()
	adminH := api.NewAdminKnowledgeHandler(kb.Ingest, kb.Knowledge)
	api.RegisterAdminKnowledgeRoutes(h, adminH, api.JWTAuthMiddleware(auth.Signer), api.AdminRoleMiddleware())

	slog.Info("listening", "addr", addr)
	h.Spin()
}
