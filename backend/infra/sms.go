package infra

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/huodaoshi/harness/backend/conf"
)

// SMSService is the interface for sending SMS messages.
type SMSService interface {
	Send(ctx context.Context, phone, code string) error
}

// LocalSMS is a no-op SMS provider that logs the message instead of sending it.
type LocalSMS struct{}

// Send logs the phone and code without performing any real SMS delivery.
func (s *LocalSMS) Send(_ context.Context, phone, code string) error {
	slog.Info("[SMS][local]", "phone", phone, "code", code)
	return nil
}

// ProdSMS is a placeholder for a real SMS provider (Aliyun / Tencent).
type ProdSMS struct {
	provider     string
	accessKey    string
	signName     string
	templateCode string
}

// Send is not yet implemented for production SMS providers.
func (s *ProdSMS) Send(_ context.Context, _, _ string) error {
	return errors.New("sms: not implemented")
}

// NewSMSService creates an SMSService based on cfg.Provider.
// "local" returns a LocalSMS that only logs.
// "aliyun" and "tencent" return a ProdSMS skeleton.
func NewSMSService(cfg conf.SMSConfig) (SMSService, error) {
	switch cfg.Provider {
	case "local":
		return &LocalSMS{}, nil
	case "aliyun", "tencent":
		return &ProdSMS{
			provider:     cfg.Provider,
			accessKey:    cfg.AccessKey,
			signName:     cfg.SignName,
			templateCode: cfg.TemplateCode,
		}, nil
	default:
		return nil, fmt.Errorf("infra: sms: unknown provider %q", cfg.Provider)
	}
}
