package session

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/huodaoshi/harness/backend/internal/safety"
)

// Executor runs the relationship session graph with SafetyGate.
type Executor struct {
	Runnable   compose.Runnable[Input, TurnOutput]
	ChatCalls  *FakeChatCallCounter
	Evaluator  *safety.Evaluator
	Templates  *safety.TemplateStore
}

// NewExecutor builds SafetyGate + branched graph (crisis | fake chat).
func NewExecutor(ctx context.Context) (*Executor, error) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		return nil, fmt.Errorf("evaluator: %w", err)
	}
	templates, err := safety.NewTemplateStore()
	if err != nil {
		return nil, fmt.Errorf("templates: %w", err)
	}
	calls := &FakeChatCallCounter{}
	runnable, err := newSessionGraph(ctx, eval, templates, calls)
	if err != nil {
		return nil, err
	}
	return &Executor{
		Runnable:  runnable,
		ChatCalls: calls,
		Evaluator: eval,
		Templates: templates,
	}, nil
}

// RunTurn executes one turn. Crisis path does not increment ChatCalls.
func (e *Executor) RunTurn(ctx context.Context, in Input) (TurnOutcome, error) {
	out, err := e.Runnable.Invoke(ctx, in)
	if err != nil {
		return TurnOutcome{}, err
	}
	if out.Crisis != nil {
		return TurnOutcome{Crisis: out.Crisis, ChatCalls: 0}, nil
	}
	return TurnOutcome{Chat: out.Chat, ChatCalls: e.ChatCalls.Load()}, nil
}

// CompileDefaultGraph keeps Spike S1 name for tests that only need streaming chat path.
func CompileDefaultGraph(ctx context.Context) (compose.Runnable[Input, string], error) {
	ex, err := NewExecutor(ctx)
	if err != nil {
		return nil, err
	}
	return newChatOnlyGraph(ctx, ex.ChatCalls)
}
