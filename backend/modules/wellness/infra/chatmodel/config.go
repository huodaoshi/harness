package chatmodel

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/huodaoshi/harness/backend/conf"
)

// Config holds LLM gateway settings (from conf + optional env overrides).
type Config struct {
	Provider           string
	FailoverProvider   string
	ARKAPIKey          string
	ARKModel           string
	ARKBaseURL         string
	RequestTimeout     time.Duration
	FirstTokenTargetMS int
}

// ConfigFromApp maps application config to gateway settings.
func ConfigFromApp(c *conf.Config) Config {
	timeout := 90 * time.Second
	if v := strings.TrimSpace(c.LLM.RequestTimeout); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			timeout = d
		}
	}
	targetMS := 3000
	if c.LLM.FirstTokenTargetMS > 0 {
		targetMS = c.LLM.FirstTokenTargetMS
	}
	provider := strings.ToLower(strings.TrimSpace(c.LLM.Provider))
	if provider == "" {
		if strings.TrimSpace(c.LLM.APIKey) != "" && strings.TrimSpace(c.LLM.Model) != "" {
			provider = "ark"
		} else {
			provider = "fake"
		}
	}
	return Config{
		Provider:           provider,
		FailoverProvider:   strings.TrimSpace(c.LLM.FailoverProvider),
		ARKAPIKey:          strings.TrimSpace(c.LLM.APIKey),
		ARKModel:           strings.TrimSpace(c.LLM.Model),
		ARKBaseURL:         strings.TrimSpace(c.LLM.BaseURL),
		RequestTimeout:     timeout,
		FirstTokenTargetMS: targetMS,
	}
}

// LoadConfigFromEnv builds config from conf.Load() or process environment (tests).
func LoadConfigFromEnv() Config {
	c, err := conf.Load()
	if err != nil {
		return legacyEnvConfig()
	}
	return ConfigFromApp(c)
}

func legacyEnvConfig() Config {
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
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// NewGatewayFromConfig builds primary (+ optional failover) gateway.
func NewGatewayFromConfig(ctx context.Context, app *conf.Config) (Gateway, Config, error) {
	cfg := ConfigFromApp(app)
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

// NewGatewayFromEnv builds gateway using conf.Load() when cwd is backend/; else process env only.
func NewGatewayFromEnv(ctx context.Context) (Gateway, Config, error) {
	c, err := conf.Load()
	if err != nil {
		cfg := legacyEnvConfig()
		primary, pErr := newProvider(ctx, cfg, cfg.Provider)
		if pErr != nil {
			return nil, cfg, pErr
		}
		return primary, cfg, nil
	}
	return NewGatewayFromConfig(ctx, c)
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
