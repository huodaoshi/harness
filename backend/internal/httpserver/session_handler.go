package httpserver

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/huodaoshi/harness/backend/internal/session"
	"github.com/huodaoshi/harness/backend/internal/store"
)

type endSessionRequest struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}

type endSessionResponse struct {
	SessionID string   `json:"session_id"`
	Summary3  []string `json:"summary3"`
}

// NewGetSessionHandler returns GET /v1/sessions/:id
func NewGetSessionHandler(st store.Store) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := requireUserID(c)
		if !ok {
			return
		}
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(consts.StatusBadRequest, map[string]string{"error": "session id required"})
			return
		}

		sess, err := st.GetSession(ctx, sessionID)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if sess == nil {
			c.JSON(consts.StatusNotFound, map[string]string{"error": "session not found"})
			return
		}
		if sess.UserID != userID {
			c.JSON(consts.StatusForbidden, map[string]string{"error": "session does not belong to user"})
			return
		}
		c.JSON(consts.StatusOK, sess)
	}
}

// NewEndSessionHandler returns POST /v1/sessions/end
func NewEndSessionHandler(st store.Store) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req endSessionRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]string{"error": "invalid json body"})
			return
		}
		userID := req.UserID
		if userID == "" {
			uid, ok := requireUserID(c)
			if !ok {
				return
			}
			userID = uid
		}
		if req.SessionID == "" {
			c.JSON(consts.StatusBadRequest, map[string]string{"error": "session_id is required"})
			return
		}

		summary3, err := session.FinalizeSession(ctx, st, req.SessionID, userID)
		if err != nil {
			writeSessionError(c, err)
			return
		}
		c.JSON(consts.StatusOK, endSessionResponse{
			SessionID: req.SessionID,
			Summary3:  summary3,
		})
	}
}

func writeSessionError(c *app.RequestContext, err error) {
	switch {
	case errors.Is(err, store.ErrSessionNotFound):
		c.JSON(consts.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, store.ErrSessionForbidden):
		c.JSON(consts.StatusForbidden, map[string]string{"error": err.Error()})
	default:
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
