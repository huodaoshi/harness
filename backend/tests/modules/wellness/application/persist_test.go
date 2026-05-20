package application_test

import (
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
	"context"
	"strings"
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/application"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
)

func TestPersist_MultiTurnThenFinalize_InjectsSummary(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	userID := "u-persist-7"

	sid, _, err := application.EnsureSession(ctx, st, "", userID, "normal")
	if err != nil {
		t.Fatal(err)
	}

	exec := newTestExecutor(t, st)

	marker := "UNIQUE_ISSUE_MARKER_7"
	_ = st.UpsertProfile(ctx, domain.RelationshipProfile{
		UserID: userID, CurrentIssue: marker,
	})

	out, err := exec.RunTurn(ctx, application.Input{UserID: userID, Message: "你好", Mode: "normal"})
	if err != nil {
		t.Fatal(err)
	}
	if err := application.PersistTurn(ctx, st, sid, userID, exec.Evaluator.Evaluate("你好"), "你好", out.Chat); err != nil {
		t.Fatal(err)
	}

	summary3, err := application.FinalizeSession(ctx, st, sid, userID)
	if err != nil || len(summary3) != 3 {
		t.Fatalf("finalize: %v %+v", err, summary3)
	}

	sid2, _, err := application.EnsureSession(ctx, st, "", userID, "normal")
	if err != nil {
		t.Fatal(err)
	}
	out2, err := exec.RunTurn(ctx, application.Input{UserID: userID, Message: "继续", Mode: "normal"})
	if err != nil {
		t.Fatal(err)
	}
	blob := out2.Chat + "\n" + out2.InjectBlock
	if !strings.Contains(blob, summary3[0]) && !strings.Contains(blob, "[上次会话摘要]") {
		t.Fatalf("expected summary inject, blob=%q summary=%v", blob, summary3)
	}
	_ = sid2
}

func TestPersist_MessageCapRejected(t *testing.T) {
	ctx := context.Background()
	st := store.NewMemoryStore()
	userID := "u-cap"

	sid, _, err := application.EnsureSession(ctx, st, "", userID, "normal")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 49; i++ {
		if err := st.AppendSessionMessages(ctx, sid, userID, "pass",
			[]domain.SessionMessage{{Role: "user", Content: "x"}}); err != nil {
			t.Fatal(err)
		}
	}

	err = st.AppendSessionMessages(ctx, sid, userID, "pass",
		[]domain.SessionMessage{{Role: "user", Content: "a"}, {Role: "assistant", Content: "b"}})
	if err != domain.ErrSessionMessageCap {
		t.Fatalf("got %v want cap", err)
	}
}
