package domain

import (
	"errors"
	"time"
)

// MaxSessionMessages is the hard cap on messages per session (PRD).
const MaxSessionMessages = 50

// SessionMessage is one turn in a conversation.
type SessionMessage struct {
	Role    string    `bson:"role" json:"role"`
	Content string    `bson:"content" json:"content"`
	At      time.Time `bson:"at" json:"at"`
}

// StoredSession is the sessions collection document (PRD shape).
type StoredSession struct {
	SessionID   string           `bson:"session_id" json:"session_id"`
	UserID      string           `bson:"user_id" json:"user_id"`
	Mode        string           `bson:"mode" json:"mode"`
	GateResults []string         `bson:"gate_results" json:"gate_results"`
	Messages    []SessionMessage `bson:"messages" json:"messages"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	EndedAt     *time.Time       `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
}

func (s *StoredSession) MessageCount() int {
	if s == nil {
		return 0
	}
	return len(s.Messages)
}

func (s *StoredSession) IsEnded() bool {
	return s != nil && s.EndedAt != nil
}

var (
	ErrSessionNotFound   = errors.New("session not found")
	ErrSessionForbidden  = errors.New("session does not belong to user")
	ErrSessionEnded      = errors.New("session already ended")
	ErrSessionMessageCap = errors.New("session message cap exceeded")
)
