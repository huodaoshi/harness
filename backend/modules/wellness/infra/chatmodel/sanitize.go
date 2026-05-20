package chatmodel

import (
	"errors"
	"strings"
)

// SanitizeError removes secrets from provider errors before sending to clients.
func SanitizeError(err error, secrets ...string) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	for _, s := range secrets {
		if s == "" {
			continue
		}
		msg = strings.ReplaceAll(msg, s, "[REDACTED]")
	}
	return errors.New(msg)
}
