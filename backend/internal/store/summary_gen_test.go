package store_test

import (
	"testing"

	"github.com/huodaoshi/harness/backend/internal/store"
)

func TestGenerateSummary3_ThreeLines(t *testing.T) {
	s := &store.StoredSession{
		Mode: "normal",
		Messages: []store.SessionMessage{
			{Role: "user", Content: "和父母吵架了"},
		},
	}
	lines := store.GenerateSummary3(s)
	if len(lines) != 3 {
		t.Fatalf("got %d lines", len(lines))
	}
}
