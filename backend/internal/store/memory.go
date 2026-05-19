package store

import (
	"context"
	"sort"
	"sync"
	"time"
)

// MemoryStore is an in-process Store for tests and offline dev without Mongo.
type MemoryStore struct {
	mu        sync.RWMutex
	profiles  map[string]RelationshipProfile
	summaries map[string][]SessionSummary // user_id -> summaries
	sessions  map[string]StoredSession
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		profiles:  make(map[string]RelationshipProfile),
		summaries: make(map[string][]SessionSummary),
		sessions:  make(map[string]StoredSession),
	}
}

func (m *MemoryStore) UpsertProfile(ctx context.Context, p RelationshipProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now().UTC()
	}
	m.profiles[p.UserID] = p
	return nil
}

func (m *MemoryStore) GetProfile(ctx context.Context, userID string) (*RelationshipProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.profiles[userID]
	if !ok {
		return nil, nil
	}
	cp := p
	return &cp, nil
}

func (m *MemoryStore) SaveSummary(ctx context.Context, s SessionSummary) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	m.summaries[s.UserID] = append(m.summaries[s.UserID], s)
	return nil
}

func (m *MemoryStore) GetLatestSummary(ctx context.Context, userID string) (*SessionSummary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := m.summaries[userID]
	if len(list) == 0 {
		return nil, nil
	}
	sorted := append([]SessionSummary(nil), list...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})
	latest := sorted[0]
	return &latest, nil
}
