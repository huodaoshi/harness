package safety_test

import (
	"testing"

	"github.com/huodaoshi/harness/backend/internal/safety"
)

func TestEvaluator_CrisisSelfHarm(t *testing.T) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		t.Fatal(err)
	}
	r := eval.Evaluate("我真的不想活了")
	if !r.IsCrisis() || r.TemplateID != "crisis_self_harm" {
		t.Fatalf("got %+v", r)
	}
}

func TestEvaluator_Pass(t *testing.T) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		t.Fatal(err)
	}
	r := eval.Evaluate("今天和父母吵了一架，很难过")
	if r.IsCrisis() {
		t.Fatalf("expected pass, got %+v", r)
	}
}
