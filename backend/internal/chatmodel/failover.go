package chatmodel

import "context"

type failoverGateway struct {
	primary      Gateway
	secondary    Gateway
	failoverName string
	secrets      []string
}

func (f *failoverGateway) Generate(ctx context.Context, req Request) (string, error) {
	text, err := f.primary.Generate(ctx, req)
	if err == nil {
		return text, nil
	}
	if f.secondary == nil {
		if f.failoverName != "" {
			return "", SanitizeError(ErrFailoverNotConfigured, f.secrets...)
		}
		return "", SanitizeError(err, f.secrets...)
	}
	text, err2 := f.secondary.Generate(ctx, req)
	if err2 != nil {
		return "", SanitizeError(err2, f.secrets...)
	}
	return text, nil
}

func (f *failoverGateway) Stream(ctx context.Context, req Request, onToken func(text string) error) error {
	err := f.primary.Stream(ctx, req, onToken)
	if err == nil {
		return nil
	}
	if f.secondary == nil {
		if f.failoverName != "" {
			return SanitizeError(ErrFailoverNotConfigured, f.secrets...)
		}
		return SanitizeError(err, f.secrets...)
	}
	return f.secondary.Stream(ctx, req, onToken)
}
