package conf

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level application configuration (P1 skeleton; grows with ADR-0002 issues).
type Config struct {
	App              AppConfig
	Log              LogConfig
	LLM              LLMConfig
	MongoDB          MongoDBConfig
	Redis            RedisConfig
	JWT              JWTConfig
	SMS              SMSConfig
	WeChat           WeChatConfig `yaml:"wechat"`
	AdminStaticLogin AdminStaticLoginConfig `yaml:"admin_static_login"`
	RateLimit        RateLimitConfig        `yaml:"rate_limit"`
}

// RateLimitConfig holds HTTP rate limit settings.
type RateLimitConfig struct {
	StreamPerMinute int `yaml:"stream_per_minute"`
}

// JWTConfig holds JWT settings.
type JWTConfig struct {
	Secret          string `yaml:"secret"`
	AccessTokenTTL  int    `yaml:"access_token_ttl"`
	RefreshTokenTTL int    `yaml:"refresh_token_ttl"`
}

// SMSConfig holds SMS provider settings.
type SMSConfig struct {
	Provider     string
	AccessKey    string `yaml:"access_key"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
}

// WeChatConfig holds WeChat mini-program settings (optional in P1).
type WeChatConfig struct {
	MiniProgram WeChatMiniProgramConfig `yaml:"miniprogram"`
}

type WeChatMiniProgramConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

// AdminStaticLoginConfig enables POST /v1/auth/admin/login with fixed credentials.
type AdminStaticLoginConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	UserID   string `yaml:"user_id"`
	UID      int64  `yaml:"uid"`
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Env  string
	Port int
}

// LogConfig holds structured logging settings (wired in a later issue).
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// LLMConfig holds LLM provider settings for wellness ChatModelGateway.
type LLMConfig struct {
	Provider string
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key"`
}

// MongoDBConfig holds MongoDB connection settings.
type MongoDBConfig struct {
	URI      string
	Database string
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Addrs    []string
	Password string
	// Required: when true, server startup pings Redis (enable for rate-limit / knowledge issues).
	Required bool `yaml:"required"`
}

// Load reads conf/config.yaml and overlays conf/{APP_ENV}.yaml (default APP_ENV=local).
func Load() (*Config, error) {
	base, err := readYAML("conf/config.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("conf: config.yaml not found (run server from backend/): %w", err)
		}
		return nil, fmt.Errorf("conf: read config.yaml: %w", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(base, cfg); err != nil {
		return nil, fmt.Errorf("conf: parse config.yaml: %w", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	cfg.App.Env = env

	overlayPath := fmt.Sprintf("conf/%s.yaml", env)
	if overlayData, oErr := readYAML(overlayPath); oErr == nil {
		overlay := &Config{}
		if err = yaml.Unmarshal(overlayData, overlay); err != nil {
			return nil, fmt.Errorf("conf: parse %s: %w", overlayPath, err)
		}
		mergeConfig(cfg, overlay)
	} else if !errors.Is(oErr, os.ErrNotExist) {
		return nil, fmt.Errorf("conf: read %s: %w", overlayPath, oErr)
	}

	if err := applyConnectionEnv(cfg); err != nil {
		return nil, err
	}
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.RateLimit.StreamPerMinute == 0 {
		cfg.RateLimit.StreamPerMinute = 60
	}
	return cfg, nil
}

func readYAML(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func mergeConfig(base, overlay *Config) {
	if overlay.App.Port != 0 {
		base.App.Port = overlay.App.Port
	}
	if overlay.App.Env != "" {
		base.App.Env = overlay.App.Env
	}
	mergeLog(&base.Log, &overlay.Log)
	mergeLLM(&base.LLM, &overlay.LLM)
	mergeMongoDB(&base.MongoDB, &overlay.MongoDB)
	mergeRedis(&base.Redis, &overlay.Redis)
	mergeJWT(&base.JWT, &overlay.JWT)
	mergeSMS(&base.SMS, &overlay.SMS)
	mergeWeChat(&base.WeChat, &overlay.WeChat)
	mergeAdminStaticLogin(&base.AdminStaticLogin, &overlay.AdminStaticLogin)
	mergeRateLimit(&base.RateLimit, &overlay.RateLimit)
}

func mergeRateLimit(base, overlay *RateLimitConfig) {
	if overlay.StreamPerMinute > 0 {
		base.StreamPerMinute = overlay.StreamPerMinute
	}
}

func mergeJWT(base, overlay *JWTConfig) {
	if overlay.Secret != "" {
		base.Secret = overlay.Secret
	}
	if overlay.AccessTokenTTL != 0 {
		base.AccessTokenTTL = overlay.AccessTokenTTL
	}
	if overlay.RefreshTokenTTL != 0 {
		base.RefreshTokenTTL = overlay.RefreshTokenTTL
	}
}

func mergeSMS(base, overlay *SMSConfig) {
	if overlay.Provider != "" {
		base.Provider = overlay.Provider
	}
	if overlay.AccessKey != "" {
		base.AccessKey = overlay.AccessKey
	}
	if overlay.SignName != "" {
		base.SignName = overlay.SignName
	}
	if overlay.TemplateCode != "" {
		base.TemplateCode = overlay.TemplateCode
	}
}

func mergeWeChat(base, overlay *WeChatConfig) {
	if overlay.MiniProgram.AppID != "" {
		base.MiniProgram.AppID = overlay.MiniProgram.AppID
	}
	if overlay.MiniProgram.AppSecret != "" {
		base.MiniProgram.AppSecret = overlay.MiniProgram.AppSecret
	}
}

func mergeAdminStaticLogin(base, overlay *AdminStaticLoginConfig) {
	if overlay.Enabled {
		base.Enabled = true
	}
	if overlay.Username != "" {
		base.Username = overlay.Username
	}
	if overlay.Password != "" {
		base.Password = overlay.Password
	}
	if overlay.UserID != "" {
		base.UserID = overlay.UserID
	}
	if overlay.UID != 0 {
		base.UID = overlay.UID
	}
}

func mergeLog(base, overlay *LogConfig) {
	if overlay.Level != "" {
		base.Level = overlay.Level
	}
	if overlay.Format != "" {
		base.Format = overlay.Format
	}
}

func mergeLLM(base, overlay *LLMConfig) {
	if overlay.Provider != "" {
		base.Provider = overlay.Provider
	}
	if overlay.Model != "" {
		base.Model = overlay.Model
	}
	if overlay.APIKey != "" {
		base.APIKey = overlay.APIKey
	}
}

func mergeMongoDB(base, overlay *MongoDBConfig) {
	if overlay.URI != "" {
		base.URI = overlay.URI
	}
	if overlay.Database != "" {
		base.Database = overlay.Database
	}
}

func mergeRedis(base, overlay *RedisConfig) {
	if len(overlay.Addrs) > 0 {
		base.Addrs = overlay.Addrs
	}
	if overlay.Password != "" {
		base.Password = overlay.Password
	}
	if overlay.Required {
		base.Required = true
	}
}
