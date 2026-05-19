package store

import "time"

// RelationshipProfile matches PRD MongoDB shape (snake_case BSON).
type RelationshipProfile struct {
	UserID        string         `bson:"user_id" json:"user_id"`
	Self          ProfileSelf    `bson:"self" json:"self"`
	People        []ProfilePerson `bson:"people" json:"people"`
	CurrentIssue  string         `bson:"current_issue" json:"current_issue"`
	UpdatedAt     time.Time      `bson:"updated_at" json:"updated_at"`
}

type ProfileSelf struct {
	Note string `bson:"note" json:"note"`
}

type ProfilePerson struct {
	Label    string `bson:"label" json:"label"`
	Relation string `bson:"relation" json:"relation"`
	Note     string `bson:"note" json:"note"`
}

// SessionSummary is Level-2 memory (latest one injected in MVP).
type SessionSummary struct {
	SessionID string    `bson:"session_id" json:"session_id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Summary3  []string  `bson:"summary3" json:"summary3"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
