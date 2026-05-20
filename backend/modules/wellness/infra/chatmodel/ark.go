package chatmodel

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

// ArkGateway calls Volcengine Ark (豆包) via eino-ext.
type ArkGateway struct {
	model    *ark.ChatModel
	timeout  time.Duration
	secrets  []string
	targetMS int
}

func NewArkGateway(ctx context.Context, cfg Config) (*ArkGateway, error) {
	if cfg.ARKAPIKey == "" || cfg.ARKModel == "" {
		return nil, errors.New("ARK_API_KEY and ARK_MODEL_ID (or LLM_MODEL_ID) are required for ark provider")
	}
	t := cfg.RequestTimeout
	httpClient := &http.Client{Timeout: t}
	chatCfg := &ark.ChatModelConfig{
		APIKey:     cfg.ARKAPIKey,
		Model:      cfg.ARKModel,
		HTTPClient: httpClient,
	}
	if cfg.ARKBaseURL != "" {
		chatCfg.BaseURL = cfg.ARKBaseURL
	}
	m, err := ark.NewChatModel(ctx, chatCfg)
	if err != nil {
		return nil, err
	}
	return &ArkGateway{
		model:    m,
		timeout:  t,
		secrets:  []string{cfg.ARKAPIKey},
		targetMS: cfg.FirstTokenTargetMS,
	}, nil
}

func (a *ArkGateway) Generate(ctx context.Context, req Request) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	msgs := buildMessages(req)
	out, err := a.model.Generate(ctx, msgs)
	if err != nil {
		return "", SanitizeError(err, a.secrets...)
	}
	if out == nil {
		return "", errors.New("empty model response")
	}
	return out.Content, nil
}

func (a *ArkGateway) Stream(ctx context.Context, req Request, onToken func(text string) error) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	msgs := buildMessages(req)
	sr, err := a.model.Stream(ctx, msgs)
	if err != nil {
		return SanitizeError(err, a.secrets...)
	}
	defer sr.Close()

	start := time.Now()
	var firstLogged bool

	for {
		msg, err := sr.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return SanitizeError(err, a.secrets...)
		}
		if msg == nil || msg.Content == "" {
			continue
		}
		if !firstLogged {
			firstLogged = true
			latency := time.Since(start).Milliseconds()
			slog.Info("llm_first_token",
				"provider", "ark",
				"latency_ms", latency,
				"target_ms", a.targetMS,
				"mode", req.Mode,
			)
		}
		for _, r := range msg.Content {
			if err := onToken(string(r)); err != nil {
				return err
			}
		}
	}
}
