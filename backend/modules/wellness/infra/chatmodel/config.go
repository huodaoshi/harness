package chatmodel

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds LLM gateway settings from environment (never commit secrets).
type Config struct {
	Provider           string
	FailoverProvider   string
	ARKAPIKey          string
	ARKModel           string
	ARKBaseURL         string
	RequestTimeout     time.Duration
	FirstTokenTargetMS int
}

func LoadConfigFromEnv() Config {
	timeout := 90 * time.Second
	if v := os.Getenv("LLM_REQUEST_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			timeout = d
		}
	}
	targetMS := 3000
	if v := os.Getenv("LLM_FIRST_TOKEN_TARGET_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			targetMS = n
		}
	}
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("LLM_PROVIDER")))
	if provider == "" {
		if os.Getenv("ARK_API_KEY") != "" && firstNonEmpty(os.Getenv("ARK_MODEL_ID"), os.Getenv("LLM_MODEL_ID")) != "" {
			provider = "ark"
		} else {
			provider = "fake"
		}
	}
	return Config{
		Provider:           provider,
		FailoverProvider:   strings.TrimSpace(os.Getenv("LLM_FAILOVER_PROVIDER")),
		ARKAPIKey:          os.Getenv("ARK_API_KEY"),
		ARKModel:           firstNonEmpty(os.Getenv("ARK_MODEL_ID"), os.Getenv("LLM_MODEL_ID")),
		ARKBaseURL:         os.Getenv("ARK_BASE_URL"),
		RequestTimeout:     timeout,
		FirstTokenTargetMS: targetMS,
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// NewGatewayFromEnv builds primary (+ optional failover) gateway.
func NewGatewayFromEnv(ctx context.Context) (Gateway, Config, error) {
	cfg := LoadConfigFromEnv()
	primary, err := newProvider(ctx, cfg, cfg.Provider)
	if err != nil {
		return nil, cfg, err
	}
	if cfg.FailoverProvider == "" {
		return primary, cfg, nil
	}
	secondary, err := newProvider(ctx, cfg, cfg.FailoverProvider)
	if err != nil {
		return &failoverGateway{
			primary:      primary,
			failoverName: cfg.FailoverProvider,
			secrets:      []string{cfg.ARKAPIKey},
		}, cfg, nil
	}
	return &failoverGateway{
		primary:   primary,
		secondary: secondary,
		secrets:   []string{cfg.ARKAPIKey},
	}, cfg, nil
}

func newProvider(ctx context.Context, cfg Config, name string) (Gateway, error) {
	switch strings.ToLower(name) {
	case "fake":
		return NewFakeGateway(), nil
	case "ark":
		return NewArkGateway(ctx, cfg)
	default:
		return nil, fmt.Errorf("%w: unknown provider %q", ErrProviderUnavailable, name)
	}
}
