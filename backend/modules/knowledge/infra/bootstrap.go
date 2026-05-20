package infra

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

const (
	// DefaultSpaceID is the P1 operational knowledge space (ADR-0002).
	DefaultSpaceID int64 = 1
	defaultSpaceName = "默认运营知识库"
)

// EnsureDefaultSpace creates space_id=1 when missing (doc_types 1/2/3).
func EnsureDefaultSpace(ctx context.Context, db *mongo.Database) (*domain.Space, error) {
	coll := db.Collection("spaces")
	var existing domain.Space
	err := coll.FindOne(ctx, bson.M{"space_id": DefaultSpaceID, "deleted_at": nil}).Decode(&existing)
	if err == nil {
		return &existing, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("infra: bootstrap space: find: %w", err)
	}

	now := time.Now().UTC()
	space := &domain.Space{
		SpaceID:     DefaultSpaceID,
		Name:        defaultSpaceName,
		Description: "P1 默认 Space：1=产品与边界，2=FAQ，3=运营话术/示例",
		DocTypes:    []int{1, 2, 3},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	counterColl := db.Collection("counters")
	if _, err := counterColl.UpdateOne(ctx, bson.M{"_id": idgen.CounterSpaceID},
		bson.M{"$max": bson.M{"seq": int64(1)}},
		options.Update().SetUpsert(true),
	); err != nil {
		return nil, fmt.Errorf("infra: bootstrap space: counter: %w", err)
	}
	if _, err := coll.InsertOne(ctx, space); err != nil {
		return nil, fmt.Errorf("infra: bootstrap space: insert: %w", err)
	}
	return space, nil
}
