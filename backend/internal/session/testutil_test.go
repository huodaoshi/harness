package session_test

import (
	"context"
	"testing"

	"github.com/huodaoshi/harness/backend/internal/chatmodel"
	"github.com/huodaoshi/harness/backend/internal/session"
	"github.com/huodaoshi/harness/backend/internal/store"
)

func newTestExecutor(t *testing.T, st store.Store) *session.Executor {
	t.Helper()
	ctx := context.Background()
	exec, err := session.NewExecutorWithGateway(ctx, st, chatmodel.NewFakeGateway(), chatmodel.Config{Provider: "fake"})
	if err != nil {
		t.Fatal(err)
	}
	return exec
}

func newTestMemoryExecutor(t *testing.T) *session.Executor {
	t.Helper()
	return newTestExecutor(t, store.NewMemoryStore())
}
