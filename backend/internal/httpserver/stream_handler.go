package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/google/uuid"

	"github.com/huodaoshi/harness/backend/internal/safety"
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

type crisisPayload struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

// NewStreamHandler returns the session SSE handler (SafetyGate + chat stream).
func NewStreamHandler(exec *session.Executor) app.HandlerFunc {
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
		exec.ChatCalls.Reset()
		outcome, err := exec.RunTurn(ctx, in)
		if err != nil {
			writeStreamError(w, err)
			return
		}

		if outcome.Crisis != nil {
			safety.Audit(sessionID, safety.Result{
				Level:      safety.Level(outcome.Crisis.TemplateID),
				TemplateID: outcome.Crisis.TemplateID,
			})
			body, err := json.Marshal(crisisPayload{
				TemplateID: outcome.Crisis.TemplateID,
				Body:       outcome.Crisis.Body,
			})
			if err != nil {
				writeStreamError(w, err)
				return
			}
			if err := w.WriteEvent("", "crisis", body); err != nil {
				fmt.Printf("sse crisis write: %v\n", err)
			}
			return
		}

		if err := session.StreamChatTokens(outcome.Chat, func(text string) error {
			body, err := json.Marshal(tokenPayload{Text: text})
			if err != nil {
				return err
			}
			return w.WriteEvent("", "token", body)
		}); err != nil {
			writeStreamError(w, err)
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

func writeStreamError(w *sse.Writer, err error) {
	errBody, _ := json.Marshal(map[string]string{
		"code":    "stream_failed",
		"message": err.Error(),
	})
	_ = w.WriteEvent("", "error", errBody)
}
