package domain

import (
	"context"
	"time"
)

// Store reads/writes relationship profiles, sessions, and session summaries.
type Store interface {
	UpsertProfile(ctx context.Context, p RelationshipProfile) error
	GetProfile(ctx context.Context, userID string) (*RelationshipProfile, error)
	SaveSummary(ctx context.Context, s SessionSummary) error
	GetLatestSummary(ctx context.Context, userID string) (*SessionSummary, error)

	GetSession(ctx context.Context, sessionID string) (*StoredSession, error)
	CreateSession(ctx context.Context, s StoredSession) error
	AppendSessionMessages(ctx context.Context, sessionID, userID string, gateResult string, msgs []SessionMessage) error
	EndSession(ctx context.Context, sessionID, userID string, endedAt time.Time) error
}
