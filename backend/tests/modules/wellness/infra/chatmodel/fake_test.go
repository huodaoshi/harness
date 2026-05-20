package chatmodel_test

import (
	"context"
	"strings"
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/infra/chatmodel"
)

func TestFakeGateway_ModeStyles(t *testing.T) {
	gw := chatmodel.NewFakeGateway()
	ctx := context.Background()

	distress, err := gw.Generate(ctx, chatmodel.Request{Mode: "distress", Message: "很累"})
	if err != nil {
		t.Fatal(err)
	}
	normal, err := gw.Generate(ctx, chatmodel.Request{Mode: "normal", Message: "很累"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(distress, "洪峰") || strings.Contains(normal, "洪峰") {
		t.Fatalf("distress=%q normal=%q", distress, normal)
	}
	if !strings.Contains(normal, "普聊") || strings.Contains(distress, "普聊") {
		t.Fatalf("distress=%q normal=%q", distress, normal)
	}
}
