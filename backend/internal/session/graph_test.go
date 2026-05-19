package session_test

import (
	"context"
	"strings"
	"testing"

	"github.com/huodaoshi/harness/backend/internal/session"
)

func TestStreamGraph_EmitsTokens(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "fake")
	ctx := context.Background()
	runnable, err := session.NewStreamGraph(ctx)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	var got strings.Builder
	err = session.StreamTurn(ctx, runnable, session.Input{
		Message: "今晚和父母吵翻了",
		Mode:    "distress",
	}, func(text string) error {
		got.WriteString(text)
		return nil
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}
	if !strings.Contains(got.String(), "洪峰") {
		t.Fatalf("unexpected reply: %q", got.String())
	}
}
