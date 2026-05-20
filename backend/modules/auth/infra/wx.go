package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/huodaoshi/harness/backend/modules/auth/domain"
)

const wxCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

// httpWeChatClient calls the WeChat jscode2session API over HTTP.
type httpWeChatClient struct {
	appID     string
	appSecret string
	client    *http.Client
}

// NewHTTPWeChatClient creates a domain.WeChatClient backed by real HTTP calls.
func NewHTTPWeChatClient(appID, appSecret string) domain.WeChatClient {
	return &httpWeChatClient{
		appID:     appID,
		appSecret: appSecret,
		client:    &http.Client{},
	}
}

// wx2SessionResponse mirrors the WeChat jscode2session JSON response.
type wx2SessionResponse struct {
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

// Code2Session exchanges a mini-program login code for openID and unionID.
func (c *httpWeChatClient) Code2Session(ctx context.Context, code string) (openID, unionID string, err error) {
	params := url.Values{}
	params.Set("appid", c.appID)
	params.Set("secret", c.appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	reqURL := wxCode2SessionURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("infra: wx: build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("infra: wx: http get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("infra: wx: read body: %w", err)
	}

	var result wx2SessionResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("infra: wx: parse response: %w", err)
	}

	if result.ErrCode != 0 {
		return "", "", fmt.Errorf("infra: wx: code2session errcode=%d errmsg=%s", result.ErrCode, result.ErrMessage)
	}

	return result.OpenID, result.UnionID, nil
}

// Compile-time assertion: httpWeChatClient must satisfy domain.WeChatClient.
var _ domain.WeChatClient = (*httpWeChatClient)(nil)
