package nextchat

import (
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/huodaoshi/harness/backend/conf"
)

// Register mounts NextChat-compatible /api/config and /api/bytedance routes.
func Register(h *server.Hertz, c *conf.Config) {
	settings := SettingsFromConfig(c)
	cfg := &ConfigHandler{Settings: settings}
	proxy := &ProxyHandler{
		Settings: settings,
		Client:   &http.Client{Timeout: 10 * time.Minute},
	}

	h.GET("/api/config", cfg.Handle)
	h.POST("/api/config", cfg.Handle)
	h.OPTIONS("/api/config", cfg.Handle)

	h.Any("/api/bytedance/*path", proxy.Handle)
}
