// Command knowledgeindexing consumes topic knowledge-ingest (RocketMQ or local in-process).
//
// Run from backend/: APP_ENV=local go run ./cmd/knowledgeindexing
//
// mq.provider=local only works when the publisher shares the same process (see cmd/server embedded worker).
// Production: mq.provider=rocketmq and run this worker separately from cmd/server.
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/knowledgebootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := conf.Load()
	if err != nil {
		log.Fatalf("conf: %v", err)
	}
	if cfg.MQ.Provider == "local" {
		slog.Warn("mq.provider=local: standalone worker cannot receive messages from another process; use cmd/server embedded worker or rocketmq")
	}

	mongoClient, _, err := infra.NewMongoClient(cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}
	defer func() { _ = mongoClient.Disconnect(context.Background()) }()

	redisClient, err := infra.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer func() { _ = redisClient.Close() }()

	cancel, err := knowledgebootstrap.StartStandaloneWorker(ctx, cfg, mongoClient, redisClient)
	if err != nil {
		log.Fatalf("worker: %v", err)
	}
	defer cancel()

	slog.Info("knowledgeindexing: subscribed", "topic", "knowledge-ingest", "mq", cfg.MQ.Provider)
	<-ctx.Done()
	slog.Info("knowledgeindexing: shutdown")
}
