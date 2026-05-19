package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

const fakeChatNode = "fake_chat"

// NewStreamGraph builds the Spike S1 Eino graph: one streamable fake chat node.
func NewStreamGraph(ctx context.Context) (compose.Runnable[Input, string], error) {
	g := compose.NewGraph[Input, string]()

	node := compose.StreamableLambda(func(ctx context.Context, in Input) (*schema.StreamReader[string], error) {
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

	if err := g.AddLambdaNode(fakeChatNode, node); err != nil {
		return nil, fmt.Errorf("add fake_chat node: %w", err)
	}
	if err := g.AddEdge(compose.START, fakeChatNode); err != nil {
		return nil, fmt.Errorf("add start edge: %w", err)
	}
	if err := g.AddEdge(fakeChatNode, compose.END); err != nil {
		return nil, fmt.Errorf("add end edge: %w", err)
	}

	return g.Compile(ctx, compose.WithGraphName("relationship_session_s1"))
}

func fakeReply(in Input) string {
	msg := strings.TrimSpace(in.Message)
	if msg == "" {
		return "我在这里，你愿意多说一点吗？"
	}
	switch in.Mode {
	case "distress":
		return "我听到了。此刻很难受是正常的，我们先慢慢说。"
	default:
		return "我听到你了：" + msg
	}
}
