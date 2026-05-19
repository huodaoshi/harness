package safety

import (
	"strings"
)

// Level is the SafetyGate outcome (L1 rules, MVP).
type Level string

const (
	LevelPass           Level = "pass"
	LevelCrisisSelfHarm Level = "crisis_self_harm"
	LevelCrisisViolence Level = "crisis_violence"
	LevelMedical        Level = "medical_boundary"
	LevelBlock          Level = "block"
)

// Result is the output of L1 SafetyGate evaluation.
type Result struct {
	Level      Level
	TemplateID string
}

func (r Result) IsCrisis() bool {
	return r.Level == LevelCrisisSelfHarm || r.Level == LevelCrisisViolence
}

func (r Result) IsMedical() bool {
	return r.Level == LevelMedical
}

func (r Result) IsBlock() bool {
	return r.Level == LevelBlock
}

// StopsLLM is true when the gate outcome must not invoke ChatModel.
func (r Result) StopsLLM() bool {
	return r.IsCrisis() || r.IsMedical() || r.IsBlock()
}

// Evaluator runs L1 keyword rules.
type Evaluator struct {
	rules []Rule
}

// NewEvaluator loads rules from the default bundled config path.
func NewEvaluator() (*Evaluator, error) {
	rules, err := LoadRules(defaultRulesPath())
	if err != nil {
		return nil, err
	}
	return &Evaluator{rules: rules}, nil
}

// Evaluate matches message text against L1 rules (first match wins, YAML order).
func (e *Evaluator) Evaluate(message string) Result {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return Result{Level: LevelPass}
	}
	for _, rule := range e.rules {
		for _, p := range rule.Patterns {
			if strings.Contains(text, strings.ToLower(p)) {
				return Result{Level: Level(rule.ID), TemplateID: rule.ID}
			}
		}
	}
	return Result{Level: LevelPass}
}
