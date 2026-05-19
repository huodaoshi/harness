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

type templatePayload struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SSE contract (#05):
//   - crisis  → { template_id, body }
//   - medical → { template_id, body }  (fixed boundary copy, no LLM)
//   - error   → { code, message }      (block + stream failures; block uses content_blocked)
//   - token / done → pass path only

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

		userID := req.UserID
		if userID == "" {
			userID = "anonymous"
		}

		in := session.Input{UserID: userID, Message: req.Message, Mode: req.Mode}
		exec.ChatCalls.Reset()
		outcome, err := exec.RunTurn(ctx, in)
		if err != nil {
			writeStreamError(w, err)
			return
		}

		gate := exec.Evaluator.Evaluate(req.Message)

		if outcome.Crisis != nil {
			safety.Audit(sessionID, gate, req.Message)
			writeTemplateEvent(w, "crisis", outcome.Crisis.TemplateID, outcome.Crisis.Body)
			return
		}

		if outcome.Medical != nil {
			safety.Audit(sessionID, gate, req.Message)
			writeTemplateEvent(w, "medical", outcome.Medical.TemplateID, outcome.Medical.Body)
			return
		}

		if outcome.Block != nil {
			safety.Audit(sessionID, gate, req.Message)
			writeGateError(w, outcome.Block.Code, outcome.Block.Message)
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

func writeTemplateEvent(w *sse.Writer, event, templateID, body string) {
	payload, err := json.Marshal(templatePayload{TemplateID: templateID, Body: body})
	if err != nil {
		writeStreamError(w, err)
		return
	}
	if err := w.WriteEvent("", event, payload); err != nil {
		fmt.Printf("sse %s write: %v\n", event, err)
	}
}

func writeGateError(w *sse.Writer, code, message string) {
	body, _ := json.Marshal(errorPayload{Code: code, Message: message})
	_ = w.WriteEvent("", "error", body)
}

func writeStreamError(w *sse.Writer, err error) {
	writeGateError(w, "stream_failed", err.Error())
}
