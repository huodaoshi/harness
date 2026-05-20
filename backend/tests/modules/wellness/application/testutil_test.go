package application_test

import (
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
	"context"
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/infra/chatmodel"
	"github.com/huodaoshi/harness/backend/modules/wellness/application"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
)

func newTestExecutor(t *testing.T, st domain.Store) *application.Executor {
	t.Helper()
	ctx := context.Background()
	exec, err := application.NewExecutorWithGateway(ctx, st, chatmodel.NewFakeGateway(), chatmodel.Config{Provider: "fake"})
	if err != nil {
		t.Fatal(err)
	}
	return exec
}

func newTestMemoryExecutor(t *testing.T) *application.Executor {
	t.Helper()
	return newTestExecutor(t, store.NewMemoryStore())
}
