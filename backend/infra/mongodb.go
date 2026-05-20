package infra

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huodaoshi/harness/backend/conf"
)

var (
	mongoOnce   sync.Once
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
	mongoErr    error
)

// NewMongoClient returns a shared MongoDB client and database handle.
func NewMongoClient(cfg conf.MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
	mongoOnce.Do(func() {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URI))
		if err != nil {
			mongoErr = fmt.Errorf("infra: mongodb connect: %w", err)
			return
		}
		if err = client.Ping(context.Background(), nil); err != nil {
			mongoErr = fmt.Errorf("infra: mongodb ping: %w", err)
			return
		}
		mongoClient = client
		mongoDB = client.Database(cfg.Database)
	})
	return mongoClient, mongoDB, mongoErr
}
