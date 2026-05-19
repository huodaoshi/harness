package httpserver

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UserIDFromRequest reads MVP tenant id from query ?user_id= or header X-User-Id.
func UserIDFromRequest(c *app.RequestContext) (string, bool) {
	if q := string(c.Query("user_id")); q != "" {
		return q, true
	}
	if h := string(c.GetHeader("X-User-Id")); h != "" {
		return h, true
	}
	return "", false
}

func requireUserID(c *app.RequestContext) (string, bool) {
	uid, ok := UserIDFromRequest(c)
	if !ok || uid == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "user_id is required (query user_id or header X-User-Id)",
		})
		return "", false
	}
	return uid, true
}
