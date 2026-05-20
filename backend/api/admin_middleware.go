package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	authdomain "github.com/huodaoshi/harness/backend/modules/auth/domain"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
	"github.com/huodaoshi/harness/backend/pkg/ctxkey"
)

// AdminRoleMiddleware requires JWT role UserRoleAdmin (from admin login).
func AdminRoleMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		v, ok := c.Get(ctxkey.Role)
		if !ok {
			apierror.Render(ctx, c, apierror.ErrForbidden)
			c.Abort()
			return
		}
		role, _ := v.(int)
		if authdomain.UserRole(role) != authdomain.UserRoleAdmin {
			apierror.Render(ctx, c, apierror.ErrForbidden)
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
