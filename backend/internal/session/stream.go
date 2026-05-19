package session

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino/compose"
)

// TokenHandler receives streamed assistant text chunks.
type TokenHandler func(text string) error

// StreamChatTokens streams a pass-path chat reply as rune-sized tokens.
func StreamChatTokens(chat string, onToken TokenHandler) error {
	for _, r := range chat {
		if err := onToken(string(r)); err != nil {
			return err
		}
	}
	return nil
}

// StreamTurn runs a streaming chat runnable (Spike S1 path).
func StreamTurn(
	ctx context.Context,
	runnable compose.Runnable[Input, string],
	in Input,
	onToken TokenHandler,
) error {
	sr, err := runnable.Stream(ctx, in)
	if err != nil {
		return err
	}
	defer sr.Close()

	for {
		chunk, err := sr.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if chunk == "" {
			continue
		}
		if err := onToken(chunk); err != nil {
			return err
		}
	}
}
