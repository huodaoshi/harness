package application

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/chatmodel"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/safety"
	"github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
)

const (
	nodeSafety        = "safety_gate"
	nodeCrisis        = "crisis_branch"
	nodeMedical       = "medical_branch"
	nodeBlock         = "block_branch"
	nodeProfileInject = "profile_inject"
	nodeFakeChat      = "fake_chat"
)

// newSessionGraph: START ? safety_gate ? branch ? crisis | medical | block | profile_inject ? chat_model ? END
func newSessionGraph(
	ctx context.Context,
	eval *safety.Evaluator,
	templates *safety.TemplateStore,
	boundary *safety.BoundaryStore,
	st domain.Store,
	gw chatmodel.Gateway,
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

	medicalNode := compose.InvokableLambda(func(ctx context.Context, routed RoutedInput) (TurnOutput, error) {
		payload, ok := boundary.RenderMedical(routed.Gate)
		if !ok {
			return TurnOutput{}, fmt.Errorf("missing medical template %q", routed.Gate.TemplateID)
		}
		return TurnOutput{Medical: &payload}, nil
	})
	if err := g.AddLambdaNode(nodeMedical, medicalNode); err != nil {
		return nil, fmt.Errorf("add medical_branch: %w", err)
	}

	blockNode := compose.InvokableLambda(func(ctx context.Context, routed RoutedInput) (TurnOutput, error) {
		payload, ok := boundary.RenderBlock(routed.Gate)
		if !ok {
			return TurnOutput{}, fmt.Errorf("missing block template %q", routed.Gate.TemplateID)
		}
		return TurnOutput{Block: &payload}, nil
	})
	if err := g.AddLambdaNode(nodeBlock, blockNode); err != nil {
		return nil, fmt.Errorf("add block_branch: %w", err)
	}

	profileNode := compose.InvokableLambda(func(ctx context.Context, routed RoutedInput) (EnrichedChatInput, error) {
		block, err := loadInjectBlock(ctx, st, routed.Input.UserID)
		if err != nil {
			return EnrichedChatInput{}, err
		}
		return EnrichedChatInput{Routed: routed, InjectBlock: block}, nil
	})
	if err := g.AddLambdaNode(nodeProfileInject, profileNode); err != nil {
		return nil, fmt.Errorf("add profile_inject: %w", err)
	}

	chatNode := compose.InvokableLambda(func(ctx context.Context, enriched EnrichedChatInput) (TurnOutput, error) {
		calls.Inc()
		req := chatmodel.Request{
			Mode:        enriched.Routed.Input.Mode,
			Message:     enriched.Routed.Input.Message,
			InjectBlock: enriched.InjectBlock,
		}
		chat, err := gw.Generate(ctx, req)
		if err != nil {
			return TurnOutput{}, err
		}
		return TurnOutput{
			Chat:        chat,
			ChatUsed:    true,
			InjectBlock: enriched.InjectBlock,
		}, nil
	})
	if err := g.AddLambdaNode(nodeFakeChat, chatNode); err != nil {
		return nil, fmt.Errorf("add chat_model: %w", err)
	}

	branch := compose.NewGraphBranch(func(ctx context.Context, routed RoutedInput) (string, error) {
		switch {
		case routed.Gate.IsCrisis():
			return nodeCrisis, nil
		case routed.Gate.IsMedical():
			return nodeMedical, nil
		case routed.Gate.IsBlock():
			return nodeBlock, nil
		default:
			return nodeProfileInject, nil
		}
	}, map[string]bool{
		nodeCrisis:        true,
		nodeMedical:       true,
		nodeBlock:         true,
		nodeProfileInject: true,
	})
	if err := g.AddBranch(nodeSafety, branch); err != nil {
		return nil, fmt.Errorf("add branch: %w", err)
	}

	if err := g.AddEdge(compose.START, nodeSafety); err != nil {
		return nil, err
	}
	if err := g.AddEdge(nodeCrisis, compose.END); err != nil {
		return nil, err
	}
	if err := g.AddEdge(nodeMedical, compose.END); err != nil {
		return nil, err
	}
	if err := g.AddEdge(nodeBlock, compose.END); err != nil {
		return nil, err
	}
	if err := g.AddEdge(nodeProfileInject, nodeFakeChat); err != nil {
		return nil, err
	}
	if err := g.AddEdge(nodeFakeChat, compose.END); err != nil {
		return nil, err
	}

	return g.Compile(ctx, compose.WithGraphName("relationship_session_s3"))
}

func loadInjectBlock(ctx context.Context, st domain.Store, userID string) (string, error) {
	if userID == "" {
		return "", nil
	}
	profile, err := st.GetProfile(ctx, userID)
	if err != nil {
		return "", err
	}
	summary, err := st.GetLatestSummary(ctx, userID)
	if err != nil {
		return "", err
	}
	return store.BuildInjectBlock(profile, summary), nil
}

// newChatOnlyGraph is the Spike S1 streaming graph (pass path only, no profile inject).
func newChatOnlyGraph(ctx context.Context, gw chatmodel.Gateway, calls *FakeChatCallCounter) (compose.Runnable[Input, string], error) {
	g := compose.NewGraph[Input, string]()

	node := compose.StreamableLambda(func(ctx context.Context, in Input) (*schema.StreamReader[string], error) {
		calls.Inc()
		text, err := gw.Generate(ctx, chatmodel.Request{Mode: in.Mode, Message: in.Message})
		if err != nil {
			return nil, err
		}
		reply := text
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

func NewStreamGraph(ctx context.Context) (compose.Runnable[Input, string], error) {
	calls := &FakeChatCallCounter{}
	gw, _, err := chatmodel.NewGatewayFromEnv(ctx)
	if err != nil {
		return nil, err
	}
	return newChatOnlyGraph(ctx, gw, calls)
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
