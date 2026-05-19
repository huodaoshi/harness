package store

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collProfiles  = "relationship_profiles"
	collSummaries = "session_summaries"
	collSessions  = "sessions"
)

// MongoStore persists profiles and summaries.
type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoStore connects using MONGODB_URI (default mongodb://localhost:27017).
func NewMongoStore(ctx context.Context) (*MongoStore, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "family_wellness"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return &MongoStore{client: client, db: client.Database(dbName)}, nil
}

func (m *MongoStore) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoStore) UpsertProfile(ctx context.Context, p RelationshipProfile) error {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now().UTC()
	}
	_, err := m.db.Collection(collProfiles).UpdateOne(
		ctx,
		bson.M{"user_id": p.UserID},
		bson.M{"$set": p},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *MongoStore) GetProfile(ctx context.Context, userID string) (*RelationshipProfile, error) {
	var p RelationshipProfile
	err := m.db.Collection(collProfiles).FindOne(ctx, bson.M{"user_id": userID}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (m *MongoStore) SaveSummary(ctx context.Context, s SessionSummary) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	_, err := m.db.Collection(collSummaries).InsertOne(ctx, s)
	return err
}

func (m *MongoStore) GetLatestSummary(ctx context.Context, userID string) (*SessionSummary, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var s SessionSummary
	err := m.db.Collection(collSummaries).FindOne(ctx, bson.M{"user_id": userID}, opts).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
