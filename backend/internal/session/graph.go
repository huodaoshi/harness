package session

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/huodaoshi/harness/backend/internal/safety"
)

const (
	nodeSafety   = "safety_gate"
	nodeCrisis   = "crisis_branch"
	nodeFakeChat = "fake_chat"
)

// newSessionGraph: START → safety_gate → branch → crisis_branch | fake_chat → END
// Crisis branch has no edge to fake_chat (ChatModel).
func newSessionGraph(
	ctx context.Context,
	eval *safety.Evaluator,
	templates *safety.TemplateStore,
	calls *FakeChatCallCounter,
) (compose.Runnable[Input, TurnOutput], error) {
	g := compose.NewGraph[Input, TurnOutput]()

	safetyNode := compose.InvokableLambda(func(ctx context.Context, in Input) (RoutedInput, error) {
		gate := eval.Evaluate(in.Message)
		return RoutedInput{Input: in, Gate: gate}, nil
	})
	if err := g.AddLambdaNode(nodeSafety, safetyNode); err != nil {
		return nil, fmt.Errorf("add safety_gate: %w", err)
	}

	crisisNode := compose.InvokableLambda(func(ctx context.Context, routed RoutedInput) (TurnOutput, error) {
		payload, ok := templates.Render(routed.Gate)
		if !ok {
			return TurnOutput{}, fmt.Errorf("missing crisis template %q", routed.Gate.TemplateID)
		}
		return TurnOutput{Crisis: &payload}, nil
	})
	if err := g.AddLambdaNode(nodeCrisis, crisisNode); err != nil {
		return nil, fmt.Errorf("add crisis_branch: %w", err)
	}

	chatNode := compose.InvokableLambda(func(ctx context.Context, routed RoutedInput) (TurnOutput, error) {
		calls.Inc()
		return TurnOutput{Chat: fakeReply(routed.Input), ChatUsed: true}, nil
	})
	if err := g.AddLambdaNode(nodeFakeChat, chatNode); err != nil {
		return nil, fmt.Errorf("add fake_chat: %w", err)
	}

	branch := compose.NewGraphBranch(func(ctx context.Context, routed RoutedInput) (string, error) {
		if routed.Gate.IsCrisis() {
			return nodeCrisis, nil
		}
		return nodeFakeChat, nil
	}, map[string]bool{nodeCrisis: true, nodeFakeChat: true})
	if err := g.AddBranch(nodeSafety, branch); err != nil {
		return nil, fmt.Errorf("add branch: %w", err)
	}

	if err := g.AddEdge(compose.START, nodeSafety); err != nil {
		return nil, fmt.Errorf("start→safety: %w", err)
	}
	if err := g.AddEdge(nodeCrisis, compose.END); err != nil {
		return nil, fmt.Errorf("crisis→end: %w", err)
	}
	if err := g.AddEdge(nodeFakeChat, compose.END); err != nil {
		return nil, fmt.Errorf("chat→end: %w", err)
	}

	return g.Compile(ctx, compose.WithGraphName("relationship_session_s2"))
}

// newChatOnlyGraph is the Spike S1 streaming graph (pass path only).
func newChatOnlyGraph(ctx context.Context, calls *FakeChatCallCounter) (compose.Runnable[Input, string], error) {
	g := compose.NewGraph[Input, string]()

	node := compose.StreamableLambda(func(ctx context.Context, in Input) (*schema.StreamReader[string], error) {
		calls.Inc()
		reply := fakeReply(in)
		sr, sw := schema.Pipe[string](len([]rune(reply)) + 1)
		go func() {
			defer sw.Close()
			for _, r := range reply {
				if closed := sw.Send(string(r), nil); closed {
					return
				}
			}
		}()
		return sr, nil
	})

	if err := g.AddLambdaNode(nodeFakeChat, node); err != nil {
		return nil, err
	}
	if err := g.AddEdge(compose.START, nodeFakeChat); err != nil {
		return nil, err
	}
	if err := g.AddEdge(nodeFakeChat, compose.END); err != nil {
		return nil, err
	}
	return g.Compile(ctx, compose.WithGraphName("relationship_session_s1_stream"))
}

// NewStreamGraph builds the Spike S1 streaming runnable (chat-only).
func NewStreamGraph(ctx context.Context) (compose.Runnable[Input, string], error) {
	calls := &FakeChatCallCounter{}
	return newChatOnlyGraph(ctx, calls)
}

func fakeReply(in Input) string {
	return fakeReplyText(in.Message, in.Mode)
}

func fakeReplyText(message, mode string) string {
	msg := trim(message)
	if msg == "" {
		return "我在这里，你愿意多说一点吗？"
	}
	switch mode {
	case "distress":
		return "我听到了。此刻很难受是正常的，我们先慢慢说。"
	default:
		return "我听到你了：" + msg
	}
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 {
		last := s[len(s)-1]
		if last != ' ' && last != '\n' && last != '\t' {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}
