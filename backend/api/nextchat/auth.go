package nextchat

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

const accessCodePrefix = "nk-"

type authResult struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg,omitempty"`
}

// Authorize applies NextChat-style access code and API key rules.
// On success, may set Authorization to the server ARK key.
func Authorize(c *app.RequestContext, s Settings) *authResult {
	authHeader := string(c.GetHeader("Authorization"))
	accessCode, apiKey := parseAuthHeader(authHeader)

	if s.NeedCode() {
		hashed := md5Hex(accessCode)
		if _, ok := s.Codes[hashed]; !ok && apiKey == "" {
			msg := "wrong access code"
			if accessCode == "" {
				msg = "empty access code"
			}
			return &authResult{Error: true, Msg: msg}
		}
	}

	if s.HideUserAPIKey() && apiKey != "" {
		return &authResult{Error: true, Msg: "you are not allowed to access with your own api key"}
	}

	if apiKey == "" && s.ARKAPIKey != "" {
		c.Request.Header.Set("Authorization", "Bearer "+s.ARKAPIKey)
	}

	return nil
}

func parseAuthHeader(bearer string) (accessCode, apiKey string) {
	token := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(bearer), "Bearer "))
	if token == "" {
		return "", ""
	}
	if strings.HasPrefix(token, accessCodePrefix) {
		return strings.TrimPrefix(token, accessCodePrefix), ""
	}
	return "", token
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
