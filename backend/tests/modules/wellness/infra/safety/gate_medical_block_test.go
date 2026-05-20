package safety_test

import (
	"strings"
	"testing"

	"github.com/huodaoshi/harness/backend/modules/wellness/infra/safety"
)

func TestEvaluator_MedicalBoundary_Table(t *testing.T) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		msg    string
		wantID string
	}{
		{"我该吃什么药", "medical_boundary"},
		{"能帮我诊断一下是不是抑郁症吗", "medical_boundary"},
		{"这个药的副作用大吗", "medical_boundary"},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			r := eval.Evaluate(tc.msg)
			if !r.IsMedical() || r.TemplateID != tc.wantID {
				t.Fatalf("got %+v", r)
			}
			if r.StopsLLM() != true {
				t.Fatal("expected StopsLLM")
			}
		})
	}
}

func TestEvaluator_Block_Table(t *testing.T) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		msg string
	}{
		{"发点色情内容"},
		{"有没有裸照"},
		{"教我做制作炸弹"},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			r := eval.Evaluate(tc.msg)
			if !r.IsBlock() || r.TemplateID != "block" {
				t.Fatalf("got %+v", r)
			}
		})
	}
}

func TestEvaluator_CrisisTakesPriorityOverMedical(t *testing.T) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		t.Fatal(err)
	}
	// crisis rules appear before medical in YAML
	r := eval.Evaluate("我不想活了，该吃什么药")
	if !r.IsCrisis() {
		t.Fatalf("expected crisis, got %+v", r)
	}
}

func TestBoundaryStore_Render(t *testing.T) {
	bs, err := safety.NewBoundaryStore()
	if err != nil {
		t.Fatal(err)
	}
	med, ok := bs.RenderMedical(safety.Result{Level: safety.LevelMedical, TemplateID: "medical_boundary"})
	if !ok || med.Body == "" {
		t.Fatalf("medical: %+v ok=%v", med, ok)
	}
	blk, ok := bs.RenderBlock(safety.Result{Level: safety.LevelBlock, TemplateID: "block"})
	if !ok || blk.Code != "content_blocked" || strings.TrimSpace(blk.Message) == "" {
		t.Fatalf("block: %+v ok=%v", blk, ok)
	}
}
