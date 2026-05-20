package chatmodel

import "errors"

var (
	ErrProviderUnavailable   = errors.New("llm provider unavailable")
	ErrFailoverNotConfigured = errors.New("failover provider is not configured")
)
