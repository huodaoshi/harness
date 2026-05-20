package nextchat

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

const apiByteDancePrefix = "/api/bytedance"

// ProxyHandler forwards /api/bytedance/* to Volcengine Ark (ARK_*).
type ProxyHandler struct {
	Settings Settings
	Client   *http.Client
}

func (h *ProxyHandler) Handle(ctx context.Context, c *app.RequestContext) {
	if string(c.Method()) == http.MethodOptions {
		c.JSON(http.StatusOK, map[string]string{"body": "OK"})
		return
	}

	if h.Settings.ARKAPIKey == "" || h.Settings.ARKModel == "" {
		c.JSON(http.StatusServiceUnavailable, map[string]any{
			"error":   true,
			"message": "ARK_API_KEY and ARK_MODEL_ID are not configured",
		})
		return
	}

	if res := Authorize(c, h.Settings); res != nil {
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	path := strings.TrimPrefix(string(c.Path()), apiByteDancePrefix)
	if path == "" {
		path = "/"
	}
	upstream := h.Settings.ARKBaseURL + path

	body, err := io.ReadAll(c.Request.BodyStream())
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	body = rewriteModelBody(body, h.Settings.ARKModel)

	req, err := http.NewRequestWithContext(ctx, string(c.Method()), upstream, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	req.Header.Set("Content-Type", string(c.ContentType()))
	if auth := string(c.GetHeader("Authorization")); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	client := h.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, map[string]any{"error": true, "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	for k, vals := range resp.Header {
		if strings.EqualFold(k, "www-authenticate") {
			continue
		}
		for _, v := range vals {
			c.Response.Header.Add(k, v)
		}
	}
	c.Response.Header.Set("X-Accel-Buffering", "no")

	if _, err := io.Copy(c.Response.BodyWriter(), resp.Body); err != nil {
		return
	}
}

func rewriteModelBody(body []byte, modelID string) []byte {
	if len(body) == 0 || modelID == "" {
		return body
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return body
	}
	if _, ok := payload["model"]; ok {
		payload["model"] = modelID
		out, err := json.Marshal(payload)
		if err != nil {
			return body
		}
		return out
	}
	return body
}
