package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/infra"
	"github.com/huodaoshi/harness/backend/infra/logging"
)

type infraBundle struct {
	Cfg         *conf.Config
	MongoClient *mongo.Client
	RedisClient redis.UniversalClient
}

func loadConfigAndInfra(ctx context.Context) (*infraBundle, error) {
	cfg, err := conf.Load()
	if err != nil {
		return nil, err
	}
	logging.Setup(cfg.Log)

	b := &infraBundle{Cfg: cfg}

	client, _, err := infra.NewMongoClient(cfg.MongoDB)
	if err != nil {
		return nil, fmt.Errorf("mongodb: %w", err)
	}
	b.MongoClient = client
	slog.Info("mongodb ok", "db", cfg.MongoDB.Database)
	if cfg.Wellness.UseMemoryStore {
		slog.Info("wellness store: in-memory", "wellness.use_memory_store", true)
	}

	if cfg.Redis.Required {
		rclient, err := infra.NewRedisClient(cfg.Redis)
		if err != nil {
			return nil, fmt.Errorf("redis: %w", err)
		}
		b.RedisClient = rclient
		slog.Info("redis ok", "addrs", cfg.Redis.Addrs)
	}

	return b, nil
}

func listenAddr(cfg *conf.Config) string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}
	return fmt.Sprintf(":%d", cfg.App.Port)
}
