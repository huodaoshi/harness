package domain

import (
	"context"
	"time"
)

// MQOutbox holds a message waiting to be published after a transient MQ failure (lazy outbox).
type MQOutbox struct {
	OutboxID string    `bson:"outbox_id"`
	JobID    string    `bson:"job_id"`
	Topic    string    `bson:"topic"`
	Payload  []byte    `bson:"payload"`
	Attempts int       `bson:"attempts"`
	LastErr  string    `bson:"last_error"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// MQOutboxRepo persists outbox rows for asynchronous MQ publish retries.
type MQOutboxRepo interface {
	// UpsertPendingAfterPublishFailure records payload for a failed Publish (job remains status=pending).
	UpsertPendingAfterPublishFailure(ctx context.Context, jobID, topic string, payload []byte) error
	ListPending(ctx context.Context, limit int64) ([]*MQOutbox, error)
	Delete(ctx context.Context, outboxID string) error
	// IncrementFailure increments attempts and sets last_error; returns new attempts value.
	IncrementFailure(ctx context.Context, outboxID string, errMsg string) (attempts int, err error)
}
