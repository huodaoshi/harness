package nextchat

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

type dangerConfig struct {
	NeedCode         bool   `json:"needCode"`
	HideUserApiKey   bool   `json:"hideUserApiKey"`
	DisableGPT4      bool   `json:"disableGPT4"`
	HideBalanceQuery bool   `json:"hideBalanceQuery"`
	DisableFastLink  bool   `json:"disableFastLink"`
	CustomModels     string `json:"customModels"`
	DefaultModel     string `json:"defaultModel"`
	VisionModels     string `json:"visionModels"`
}

// ConfigHandler serves GET/POST /api/config (NextChat-compatible).
type ConfigHandler struct {
	Settings Settings
}

func (h *ConfigHandler) Handle(ctx context.Context, c *app.RequestContext) {
	body := dangerConfig{
		NeedCode:         h.Settings.NeedCode(),
		HideUserApiKey:   h.Settings.HideUserAPIKey(),
		DisableGPT4:      false,
		HideBalanceQuery: true,
		DisableFastLink:  false,
		CustomModels:     h.Settings.CustomModels,
		DefaultModel:     h.Settings.DefaultModel,
		VisionModels:     "",
	}
	c.JSON(http.StatusOK, body)
}
