package store

import (
	"context"
	"errors"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStore) GetSession(ctx context.Context, sessionID string) (*domain.StoredSession, error) {
	var s domain.StoredSession
	err := m.db.Collection(collSessions).FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *MongoStore) CreateSession(ctx context.Context, s domain.StoredSession) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	if s.GateResults == nil {
		s.GateResults = []string{}
	}
	if s.Messages == nil {
		s.Messages = []domain.SessionMessage{}
	}
	_, err := m.db.Collection(collSessions).InsertOne(ctx, s)
	return err
}

func (m *MongoStore) AppendSessionMessages(ctx context.Context, sessionID, userID string, gateResult string, msgs []domain.SessionMessage) error {
	sess, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if sess == nil {
		return domain.ErrSessionNotFound
	}
	if sess.UserID != userID {
		return domain.ErrSessionForbidden
	}
	if sess.IsEnded() {
		return domain.ErrSessionEnded
	}
	if len(sess.Messages)+len(msgs) > domain.MaxSessionMessages {
		return domain.ErrSessionMessageCap
	}

	now := time.Now().UTC()
	setMsgs := make([]domain.SessionMessage, len(msgs))
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
		return domain.ErrSessionNotFound
	}
	return nil
}
