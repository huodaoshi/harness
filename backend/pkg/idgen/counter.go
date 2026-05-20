package idgen

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Counter name constants for well-known sequence namespaces.
const (
	CounterUID         = "uid"
	CounterSpaceID     = "space_id"
	CounterAdminDeptID = "admin_dept_id"
	CounterAdminRoleID = "admin_role_id"
)

// counterDoc is the MongoDB document shape for the counters collection.
type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int64  `bson:"seq"`
}

// MongoCounter is a MongoDB-backed auto-increment sequence generator.
type MongoCounter struct {
	coll *mongo.Collection
}

// NewMongoCounter creates a MongoCounter backed by the "counters" collection in db.
func NewMongoCounter(db *mongo.Database) *MongoCounter {
	return &MongoCounter{coll: db.Collection("counters")}
}

// Next atomically increments the named counter and returns the new value.
func (c *MongoCounter) Next(ctx context.Context, name string) (int64, error) {
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"seq": int64(1)}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var doc counterDoc
	if err := c.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return 0, fmt.Errorf("idgen: counter: next %q failed: %w", name, err)
	}
	return doc.Seq, nil
}
