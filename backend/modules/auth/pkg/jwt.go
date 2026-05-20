package authpkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSigner issues and parses JWT tokens.
type JWTSigner interface {
	Sign(userID string, uid int64, role int) (string, error)
	SignWithTTL(userID string, uid int64, role int, ttl time.Duration) (string, error)
	Parse(token string) (userID string, uid int64, role int, err error)
}

// customClaims embeds standard JWT claims and adds application-specific fields.
type customClaims struct {
	UID  int64 `json:"uid"`
	Role int   `json:"role"`
	jwt.RegisteredClaims
}

// hs256Signer is an HMAC-SHA256 JWT signer.
type hs256Signer struct {
	secret         []byte
	accessTokenTTL time.Duration
}

// NewHS256Signer creates a JWTSigner using HMAC-SHA256.
// accessTokenTTL is in seconds.
func NewHS256Signer(secret string, accessTokenTTL int) JWTSigner {
	return &hs256Signer{
		secret:         []byte(secret),
		accessTokenTTL: time.Duration(accessTokenTTL) * time.Second,
	}
}

// Sign creates a signed JWT token for the given user using the configured access token TTL.
func (s *hs256Signer) Sign(userID string, uid int64, role int) (string, error) {
	return s.SignWithTTL(userID, uid, role, s.accessTokenTTL)
}

// SignWithTTL creates a signed JWT token for the given user with an explicit TTL.
func (s *hs256Signer) SignWithTTL(userID string, uid int64, role int, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := customClaims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("authpkg: jwt: sign failed: %w", err)
	}
	return signed, nil
}

// Parse validates and parses a JWT token, returning its claims.
func (s *hs256Signer) Parse(tokenStr string) (userID string, uid int64, role int, err error) {
	token, parseErr := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("authpkg: jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if parseErr != nil {
		return "", 0, 0, fmt.Errorf("authpkg: jwt: parse failed: %w", parseErr)
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return "", 0, 0, fmt.Errorf("authpkg: jwt: invalid claims")
	}

	return claims.Subject, claims.UID, claims.Role, nil
}
