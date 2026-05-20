package nextchat

import (
	"fmt"
	"strings"

	"github.com/huodaoshi/harness/backend/conf"
)

const (
	defaultARKBaseURL = "https://ark.cn-beijing.volces.com"
	defaultModelLabel = "Harness"
)

// Settings holds NextChat-compatible server-side options.
type Settings struct {
	ARKAPIKey    string
	ARKModel     string
	ARKBaseURL   string
	Codes        map[string]struct{}
	CustomModels string
	DefaultModel string
}

// SettingsFromConfig builds proxy settings from loaded application config.
func SettingsFromConfig(c *conf.Config) Settings {
	key := strings.TrimSpace(c.LLM.APIKey)
	model := strings.TrimSpace(c.LLM.Model)
	base := strings.TrimSpace(c.LLM.BaseURL)
	if base == "" {
		base = defaultARKBaseURL
	}

	custom := strings.TrimSpace(c.NextChat.CustomModels)
	defaultModel := strings.TrimSpace(c.NextChat.DefaultModel)
	if custom == "" && model != "" {
		label := defaultModelLabel
		if defaultModel != "" {
			if before, _, ok := strings.Cut(defaultModel, "@"); ok {
				label = before
			} else {
				label = defaultModel
			}
		}
		custom = fmt.Sprintf("-all,+%s@bytedance=%s", label, model)
		defaultModel = fmt.Sprintf("%s@bytedance", label)
	} else if defaultModel == "" && model != "" {
		defaultModel = fmt.Sprintf("%s@bytedance", defaultModelLabel)
	}

	codes := make(map[string]struct{})
	for _, part := range c.NextChat.AccessCodes {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		codes[md5Hex(part)] = struct{}{}
	}

	return Settings{
		ARKAPIKey:    key,
		ARKModel:     model,
		ARKBaseURL:   strings.TrimRight(base, "/"),
		Codes:        codes,
		CustomModels: custom,
		DefaultModel: defaultModel,
	}
}

func (s Settings) NeedCode() bool {
	return len(s.Codes) > 0
}

func (s Settings) HideUserAPIKey() bool {
	return s.ARKAPIKey != ""
}
