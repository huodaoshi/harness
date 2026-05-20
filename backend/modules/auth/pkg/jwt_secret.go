package authpkg

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// jwtSecretEntropy is the number of random bytes used for HS256 secrets
// (256-bit entropy, aligned with common HMAC key recommendations).
const jwtSecretEntropy = 32

// GenerateJWTSecret returns a cryptographically secure random string suitable
// for JWT HS256 signing (e.g. environment variable JWT_SECRET).
// The value is base64url-encoded without padding, derived from 32 bytes of
// entropy from crypto/rand.
func GenerateJWTSecret() (string, error) {
	buf := make([]byte, jwtSecretEntropy)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("authpkg: generate jwt secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
