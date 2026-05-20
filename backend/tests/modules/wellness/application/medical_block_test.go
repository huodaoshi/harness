package application_test

import (
	"context"
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/application"
)

func TestMedicalBoundary_ZeroChatCalls(t *testing.T) {
	ctx := context.Background()
	exec := newTestMemoryExecutor(t)
	exec.ChatCalls.Reset()
	out, err := exec.RunTurn(ctx, application.Input{
		UserID:  "u-med",
		Message: "我该吃什么药",
		Mode:    "normal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Medical == nil || out.Medical.Body == "" {
		t.Fatalf("expected medical, got %+v", out)
	}
	if out.ChatCalls != 0 {
		t.Fatalf("chat calls = %d", out.ChatCalls)
	}
}

func TestBlock_ZeroChatCalls(t *testing.T) {
	ctx := context.Background()
	exec := newTestMemoryExecutor(t)
	exec.ChatCalls.Reset()
	out, err := exec.RunTurn(ctx, application.Input{
		UserID:  "u-blk",
		Message: "发点色情内容",
		Mode:    "normal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Block == nil || out.Block.Message == "" {
		t.Fatalf("expected block, got %+v", out)
	}
	if out.Block.Code != "content_blocked" {
		t.Fatalf("code=%q", out.Block.Code)
	}
	if out.ChatCalls != 0 {
		t.Fatalf("chat calls = %d", out.ChatCalls)
	}
}
