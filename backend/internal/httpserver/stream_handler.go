package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/google/uuid"

	"github.com/huodaoshi/harness/backend/internal/session"
)

// StreamRequest is the JSON body for POST /v1/sessions/stream.
type StreamRequest struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Mode      string `json:"mode"`
	Message   string `json:"message"`
}

type tokenPayload struct {
	Text string `json:"text"`
}

type donePayload struct {
	SessionID string `json:"session_id"`
}

// NewStreamHandler returns the Spike S1 SSE handler.
func NewStreamHandler(runnable compose.Runnable[session.Input, string]) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req StreamRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]string{"error": "invalid json body"})
			return
		}
		if req.Message == "" {
			c.JSON(consts.StatusBadRequest, map[string]string{"error": "message is required"})
			return
		}
		if req.Mode == "" {
			req.Mode = "normal"
		}

		sessionID := req.SessionID
		if sessionID == "" {
			sessionID = uuid.NewString()
		}

		c.SetStatusCode(http.StatusOK)
		w := sse.NewWriter(c)

		in := session.Input{Message: req.Message, Mode: req.Mode}
		err := session.StreamTurn(ctx, runnable, in, func(text string) error {
			body, err := json.Marshal(tokenPayload{Text: text})
			if err != nil {
				return err
			}
			return w.WriteEvent("", "token", body)
		})
		if err != nil {
			errBody, _ := json.Marshal(map[string]string{
				"code":    "stream_failed",
				"message": err.Error(),
			})
			_ = w.WriteEvent("", "error", errBody)
			return
		}

		done, err := json.Marshal(donePayload{SessionID: sessionID})
		if err != nil {
			return
		}
		if err := w.WriteEvent("", "done", done); err != nil {
			fmt.Printf("sse done write: %v\n", err)
		}
	}
}
