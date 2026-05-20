package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/safety"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
)

// EnsureSession loads an active session or creates one. Returns session id to use.
func EnsureSession(ctx context.Context, st domain.Store, sessionID, userID, mode string) (string, *domain.StoredSession, error) {
	if sessionID != "" {
		sess, err := st.GetSession(ctx, sessionID)
		if err != nil {
			return "", nil, err
		}
		if sess != nil {
			if sess.UserID != userID {
				return "", nil, domain.ErrSessionForbidden
			}
			if sess.IsEnded() {
				return "", nil, domain.ErrSessionEnded
			}
			return sess.SessionID, sess, nil
		}
	}

	id := sessionID
	if id == "" {
		id = uuid.NewString()
	}
	sess := domain.StoredSession{
		SessionID:   id,
		UserID:      userID,
		Mode:        mode,
		GateResults: []string{},
		Messages:    []domain.SessionMessage{},
		CreatedAt:   time.Now().UTC(),
	}
	if err := st.CreateSession(ctx, sess); err != nil {
		return "", nil, err
	}
	return id, &sess, nil
}

// PersistTurn appends user + assistant messages for one stream turn.
func PersistTurn(
	ctx context.Context,
	st domain.Store,
	sessionID, userID string,
	gate safety.Result,
	userText, assistantText string,
) error {
	gateResult := string(gate.Level)
	msgs := []domain.SessionMessage{
		{Role: "user", Content: userText},
	}
	if assistantText != "" {
		msgs = append(msgs, domain.SessionMessage{Role: "assistant", Content: assistantText})
	}
	return st.AppendSessionMessages(ctx, sessionID, userID, gateResult, msgs)
}

// FinalizeSession ends the session and writes summary3 (PostProcess MVP).
func FinalizeSession(ctx context.Context, st domain.Store, sessionID, userID string) ([]string, error) {
	sess, err := st.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, domain.ErrSessionNotFound
	}
	if sess.UserID != userID {
		return nil, domain.ErrSessionForbidden
	}

	now := time.Now().UTC()
	if !sess.IsEnded() {
		if err := st.EndSession(ctx, sessionID, userID, now); err != nil {
			return nil, err
		}
		sess, err = st.GetSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
	}

	existing, err := st.GetLatestSummary(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.SessionID == sessionID {
		return existing.Summary3, nil
	}

	summary3 := store.GenerateSummary3(sess)
	sum := domain.SessionSummary{
		SessionID: sessionID,
		UserID:    userID,
		Summary3:  summary3,
		CreatedAt: now,
	}
	if err := st.SaveSummary(ctx, sum); err != nil {
		return nil, err
	}
	return summary3, nil
}
