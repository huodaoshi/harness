package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	authpkg "github.com/huodaoshi/harness/backend/modules/auth/pkg"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
	"github.com/huodaoshi/harness/backend/pkg/ctxkey"
)

// JWTOrGuestMiddleware accepts either a valid Bearer JWT or a valid X-Anon-ID (guest).
func JWTOrGuestMiddleware(signer authpkg.JWTSigner) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				apierror.Render(ctx, c, apierror.ErrTokenExpired)
				c.Abort()
				return
			}
			userID, uid, role, err := signer.Parse(parts[1])
			if err != nil {
				apierror.Render(ctx, c, apierror.ErrTokenExpired)
				c.Abort()
				return
			}
			c.Set(ctxkey.UserID, userID)
			c.Set(ctxkey.UID, uid)
			c.Set(ctxkey.Role, role)
			c.Set(ctxkey.ClientIP, extractIP(c))
			platform := strings.TrimSpace(string(c.GetHeader("X-Platform")))
			if platform == "" {
				platform = "h5"
			}
			c.Set(ctxkey.Platform, platform)
			c.Next(ctx)
			return
		}

		anonID := strings.TrimSpace(string(c.GetHeader("X-Anon-ID")))
		if anonID == "" {
			apierror.Render(ctx, c, apierror.ErrTokenExpired)
			c.Abort()
			return
		}
		if !isValidAnonUUID(anonID) {
			apierror.Render(ctx, c, apierror.ErrAnonIDInvalid)
			c.Abort()
			return
		}
		c.Set(ctxkey.UserID, "anon:"+anonID)
		c.Set(ctxkey.UID, int64(0))
		c.Set(ctxkey.Role, 0)
		c.Set(ctxkey.AnonID, anonID)
		c.Next(ctx)
	}
}

// JWTAuthMiddleware validates the Bearer token and writes user claims into the context.
func JWTAuthMiddleware(signer authpkg.JWTSigner) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			apierror.Render(ctx, c, apierror.ErrTokenExpired)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			apierror.Render(ctx, c, apierror.ErrTokenExpired)
			c.Abort()
			return
		}

		userID, uid, role, err := signer.Parse(parts[1])
		if err != nil {
			apierror.Render(ctx, c, apierror.ErrTokenExpired)
			c.Abort()
			return
		}

		c.Set(ctxkey.UserID, userID)
		c.Set(ctxkey.UID, uid)
		c.Set(ctxkey.Role, role)
		c.Set(ctxkey.ClientIP, extractIP(c))
		platform := strings.TrimSpace(string(c.GetHeader("X-Platform")))
		if platform == "" {
			platform = "h5"
		}
		c.Set(ctxkey.Platform, platform)

		c.Next(ctx)
	}
}

func extractUserID(c *app.RequestContext) string {
	v, _ := c.Get(ctxkey.UserID)
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func extractUID(c *app.RequestContext) int64 {
	v, _ := c.Get(ctxkey.UID)
	if v == nil {
		return 0
	}
	uid, _ := v.(int64)
	return uid
}

func extractAnonID(c *app.RequestContext) string {
	v, _ := c.Get(ctxkey.AnonID)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
