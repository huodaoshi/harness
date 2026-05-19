package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStore) GetSession(ctx context.Context, sessionID string) (*StoredSession, error) {
	var s StoredSession
	err := m.db.Collection(collSessions).FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *MongoStore) CreateSession(ctx context.Context, s StoredSession) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	if s.GateResults == nil {
		s.GateResults = []string{}
	}
	if s.Messages == nil {
		s.Messages = []SessionMessage{}
	}
	_, err := m.db.Collection(collSessions).InsertOne(ctx, s)
	return err
}

func (m *MongoStore) AppendSessionMessages(ctx context.Context, sessionID, userID string, gateResult string, msgs []SessionMessage) error {
	sess, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if sess == nil {
		return ErrSessionNotFound
	}
	if sess.UserID != userID {
		return ErrSessionForbidden
	}
	if sess.IsEnded() {
		return ErrSessionEnded
	}
	if len(sess.Messages)+len(msgs) > MaxSessionMessages {
		return ErrSessionMessageCap
	}

	now := time.Now().UTC()
	setMsgs := make([]SessionMessage, len(msgs))
	copy(setMsgs, msgs)
	for i := range setMsgs {
		if setMsgs[i].At.IsZero() {
			setMsgs[i].At = now
		}
	}

	pushDoc := bson.M{
		"messages": bson.M{"$each": setMsgs},
	}
	if gateResult != "" {
		pushDoc["gate_results"] = gateResult
	}

	_, err = m.db.Collection(collSessions).UpdateOne(
		ctx,
		bson.M{"session_id": sessionID, "user_id": userID, "ended_at": nil},
		bson.M{"$push": pushDoc},
	)
	return err
}

func (m *MongoStore) EndSession(ctx context.Context, sessionID, userID string, endedAt time.Time) error {
	t := endedAt.UTC()
	res, err := m.db.Collection(collSessions).UpdateOne(
		ctx,
		bson.M{"session_id": sessionID, "user_id": userID},
		bson.M{"$set": bson.M{"ended_at": t}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrSessionNotFound
	}
	return nil
}
