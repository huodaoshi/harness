package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/infra"
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

	b := &infraBundle{Cfg: cfg}

	client, _, err := infra.NewMongoClient(cfg.MongoDB)
	if err != nil {
		return nil, fmt.Errorf("mongodb: %w", err)
	}
	b.MongoClient = client
	log.Printf("mongodb ok (db=%s)", cfg.MongoDB.Database)
	if os.Getenv("USE_MEMORY_STORE") == "true" {
		log.Printf("wellness store: in-memory (USE_MEMORY_STORE=true)")
	}

	if cfg.Redis.Required {
		rclient, err := infra.NewRedisClient(cfg.Redis)
		if err != nil {
			return nil, fmt.Errorf("redis: %w", err)
		}
		b.RedisClient = rclient
		log.Printf("redis ok (addrs=%v)", cfg.Redis.Addrs)
	}

	return b, nil
}

func listenAddr(cfg *conf.Config) string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}
	return fmt.Sprintf(":%d", cfg.App.Port)
}
