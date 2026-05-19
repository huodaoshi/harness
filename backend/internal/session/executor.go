package session

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/huodaoshi/harness/backend/internal/safety"
	"github.com/huodaoshi/harness/backend/internal/store"
)

// Executor runs the relationship session graph with SafetyGate and ProfileInject.
type Executor struct {
	Runnable  compose.Runnable[Input, TurnOutput]
	ChatCalls *FakeChatCallCounter
	Evaluator *safety.Evaluator
	Templates *safety.TemplateStore
	Boundary  *safety.BoundaryStore
	Store     store.Store
}

// NewExecutor builds the full S3 graph with store from environment.
func NewExecutor(ctx context.Context) (*Executor, error) {
	st, err := store.NewFromEnv(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return NewExecutorWithStore(ctx, st)
}

// NewExecutorWithStore is for tests and explicit wiring.
func NewExecutorWithStore(ctx context.Context, st store.Store) (*Executor, error) {
	eval, err := safety.NewEvaluator()
	if err != nil {
		return nil, fmt.Errorf("evaluator: %w", err)
	}
	templates, err := safety.NewTemplateStore()
	if err != nil {
		return nil, fmt.Errorf("templates: %w", err)
	}
	boundary, err := safety.NewBoundaryStore()
	if err != nil {
		return nil, fmt.Errorf("boundary: %w", err)
	}
	calls := &FakeChatCallCounter{}
	runnable, err := newSessionGraph(ctx, eval, templates, boundary, st, calls)
	if err != nil {
		return nil, err
	}
	return &Executor{
		Runnable:  runnable,
		ChatCalls: calls,
		Evaluator: eval,
		Templates: templates,
		Boundary:  boundary,
		Store:     st,
	}, nil
}

// RunTurn executes one turn. Gate branches do not increment ChatCalls.
func (e *Executor) RunTurn(ctx context.Context, in Input) (TurnOutcome, error) {
	out, err := e.Runnable.Invoke(ctx, in)
	if err != nil {
		return TurnOutcome{}, err
	}
	if out.Crisis != nil {
		return TurnOutcome{Crisis: out.Crisis, ChatCalls: 0}, nil
	}
	if out.Medical != nil {
		return TurnOutcome{Medical: out.Medical, ChatCalls: 0}, nil
	}
	if out.Block != nil {
		return TurnOutcome{Block: out.Block, ChatCalls: 0}, nil
	}
	return TurnOutcome{
		Chat:        out.Chat,
		ChatCalls:   e.ChatCalls.Load(),
		InjectBlock: out.InjectBlock,
	}, nil
}

func CompileDefaultGraph(ctx context.Context) (compose.Runnable[Input, string], error) {
	calls := &FakeChatCallCounter{}
	return newChatOnlyGraph(ctx, calls)
}
