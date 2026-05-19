package chatmodel

import "context"

// Gateway streams or generates assistant text for the pass path (never crisis/medical/block).
type Gateway interface {
	Generate(ctx context.Context, req Request) (string, error)
	Stream(ctx context.Context, req Request, onToken func(text string) error) error
}
