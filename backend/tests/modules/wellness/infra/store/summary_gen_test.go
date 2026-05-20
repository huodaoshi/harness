package store_test

import (
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
)

func TestGenerateSummary3_ThreeLines(t *testing.T) {
	s := &domain.StoredSession{
		Mode: "normal",
		Messages: []domain.SessionMessage{
			{Role: "user", Content: "和父母吵架了"},
		},
	}
	lines := store.GenerateSummary3(s)
	if len(lines) != 3 {
		t.Fatalf("got %d lines", len(lines))
	}
}
