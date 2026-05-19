package store

import (
	"context"
	"time"
)

func (m *MemoryStore) GetSession(ctx context.Context, sessionID string) (*StoredSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	cp := cloneSession(s)
	return &cp, nil
}

func (m *MemoryStore) CreateSession(ctx context.Context, s StoredSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	if s.GateResults == nil {
		s.GateResults = []string{}
	}
	if s.Messages == nil {
		s.Messages = []SessionMessage{}
	}
	m.sessions[s.SessionID] = s
	return nil
}

func (m *MemoryStore) AppendSessionMessages(ctx context.Context, sessionID, userID string, gateResult string, msgs []SessionMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}
	if s.UserID != userID {
		return ErrSessionForbidden
	}
	if s.EndedAt != nil {
		return ErrSessionEnded
	}
	if len(s.Messages)+len(msgs) > MaxSessionMessages {
		return ErrSessionMessageCap
	}
	now := time.Now().UTC()
	for i := range msgs {
		if msgs[i].At.IsZero() {
			msgs[i].At = now
		}
	}
	s.Messages = append(s.Messages, msgs...)
	if gateResult != "" {
		s.GateResults = append(s.GateResults, gateResult)
	}
	m.sessions[sessionID] = s
	return nil
}

func (m *MemoryStore) EndSession(ctx context.Context, sessionID, userID string, endedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}
	if s.UserID != userID {
		return ErrSessionForbidden
	}
	if s.EndedAt != nil {
		return nil
	}
	t := endedAt.UTC()
	s.EndedAt = &t
	m.sessions[sessionID] = s
	return nil
}

func cloneSession(s StoredSession) StoredSession {
	cp := s
	if s.GateResults != nil {
		cp.GateResults = append([]string(nil), s.GateResults...)
	}
	if s.Messages != nil {
		cp.Messages = append([]SessionMessage(nil), s.Messages...)
	}
	return cp
}
