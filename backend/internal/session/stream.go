package session

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino/compose"
)

// TokenHandler receives streamed assistant text chunks.
type TokenHandler func(text string) error

// StreamTurn runs the compiled graph and invokes handler for each token chunk.
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

// CompileDefaultGraph compiles the Spike S1 graph for production wiring.
func CompileDefaultGraph(ctx context.Context) (compose.Runnable[Input, string], error) {
	return NewStreamGraph(ctx)
}
