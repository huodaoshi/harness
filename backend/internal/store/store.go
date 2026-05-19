package store

import "context"

// Store reads/writes relationship profiles and session summaries.
type Store interface {
	UpsertProfile(ctx context.Context, p RelationshipProfile) error
	GetProfile(ctx context.Context, userID string) (*RelationshipProfile, error)
	SaveSummary(ctx context.Context, s SessionSummary) error
	GetLatestSummary(ctx context.Context, userID string) (*SessionSummary, error)
}
