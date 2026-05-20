package api

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
	"github.com/huodaoshi/harness/backend/pkg/ctxkey"
)

// UserIDFromContext returns the authenticated tenant id set by JWTOrGuestMiddleware.
func UserIDFromContext(c *app.RequestContext) (string, bool) {
	v, ok := c.Get(ctxkey.UserID)
	if !ok || v == nil {
		return "", false
	}
	uid := fmt.Sprintf("%v", v)
	return uid, uid != ""
}

func requireUserID(c *app.RequestContext) (string, bool) {
	uid, ok := UserIDFromContext(c)
	if !ok {
		apierror.Render(nil, c, apierror.ErrTokenExpired)
		return "", false
	}
	return uid, true
}

