package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/huodaoshi/harness/backend/internal/chatmodel"
)

// StreamPassTurn runs pass-path ProfileInject + ChatModel streaming (no gate branches).
func (e *Executor) StreamPassTurn(
	ctx context.Context,
	in Input,
	onToken func(text string) error,
) (fullText string, injectBlock string, err error) {
	gate := e.Evaluator.Evaluate(in.Message)
	if gate.StopsLLM() {
		return "", "", fmt.Errorf("stream pass called on gate branch %s", gate.Level)
	}

	injectBlock, err = loadInjectBlock(ctx, e.Store, in.UserID)
	if err != nil {
		return "", "", err
	}

	e.ChatCalls.Inc()
	req := chatmodel.Request{
		Mode:        in.Mode,
		Message:     in.Message,
		InjectBlock: injectBlock,
	}

	var b strings.Builder
	err = e.Gateway.Stream(ctx, req, func(tok string) error {
		b.WriteString(tok)
		return onToken(tok)
	})
	return b.String(), injectBlock, err
}
