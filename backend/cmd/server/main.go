package main

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/api"
	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/knowledgebootstrap"
	"github.com/huodaoshi/harness/backend/modules/wellness/application"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/configpaths"
)

func main() {
	ctx := context.Background()

	bundle, err := loadConfigAndInfra(ctx)
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	log.Printf("config loaded (env=%s port=%d)", bundle.Cfg.App.Env, bundle.Cfg.App.Port)

	auth, err := wireAuth(bundle.Cfg, bundle.MongoClient, bundle.RedisClient)
	if err != nil {
		log.Fatalf("auth: %v", err)
	}

	exec, err := application.NewExecutor(ctx)
	if err != nil {
		log.Fatalf("session executor: %v", err)
	}

	addr := listenAddr(bundle.Cfg)

	h := server.Default(server.WithHostPorts(addr))
	api.RegisterAuthRoutes(h, api.NewAuthHandler(auth.Service), api.JWTAuthMiddleware(auth.Signer))
	streamRL := infra.NewRedisRateLimiter(bundle.RedisClient)
	api.RegisterWellnessRoutes(h, exec, api.JWTOrGuestMiddleware(auth.Signer), streamRL, bundle.Cfg.RateLimit.StreamPerMinute)

	kb, err := knowledgebootstrap.Wire(ctx, bundle.Cfg, bundle.MongoClient, bundle.RedisClient, true)
	if err != nil {
		log.Fatalf("knowledge: %v", err)
	}
	defer kb.Close()
	adminH := api.NewAdminKnowledgeHandler(kb.Ingest, kb.Knowledge)
	api.RegisterAdminKnowledgeRoutes(h, adminH, api.JWTAuthMiddleware(auth.Signer), api.AdminRoleMiddleware())

	log.Printf("listening on %s (web=%s)", addr, configpaths.WebRoot())
	h.Spin()
}
