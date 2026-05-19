package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/sse"

	"github.com/huodaoshi/harness/backend/internal/chatmodel"
	"github.com/huodaoshi/harness/backend/internal/safety"
	"github.com/huodaoshi/harness/backend/internal/session"
	"github.com/huodaoshi/harness/backend/internal/store"
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
	SessionID    string `json:"session_id"`
	MessageCount int    `json:"message_count"`
}

type templatePayload struct {
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

		userID := req.UserID
		if userID == "" {
			userID = "anonymous"
		}

		st := exec.Store
		sessionID, sess, err := session.EnsureSession(ctx, st, req.SessionID, userID, req.Mode)
		if err != nil {
			writeSessionHTTPError(c, err)
			return
		}

		incoming := 2
		if sess != nil && sess.MessageCount()+incoming > store.MaxSessionMessages {
			c.JSON(consts.StatusConflict, map[string]string{
				"error": "session message cap exceeded",
				"code":  "session_message_cap",
			})
			return
		}

		c.SetStatusCode(http.StatusOK)
		w := sse.NewWriter(c)

		in := session.Input{UserID: userID, Message: req.Message, Mode: req.Mode}
		gate := exec.Evaluator.Evaluate(req.Message)
		exec.ChatCalls.Reset()

		if gate.StopsLLM() {
			handleGateTurn(ctx, w, exec, st, sessionID, userID, in, gate, req.Message)
			return
		}

		fullChat, _, err := exec.StreamPassTurn(ctx, in, func(text string) error {
			body, err := json.Marshal(tokenPayload{Text: text})
			if err != nil {
				return err
			}
			return w.WriteEvent("", "token", body)
		})
		if err != nil {
			writeProviderError(w, exec, err)
			return
		}

		if err := session.PersistTurn(ctx, st, sessionID, userID, gate, req.Message, fullChat); err != nil {
			if errors.Is(err, store.ErrSessionMessageCap) {
				writeGateError(w, "session_message_cap", "本场对话已达到消息上限（50 条），请结束本次对话后开始新会话。")
				return
			}
			writeProviderError(w, exec, err)
			return
		}

		writeDone(w, sessionID, st, ctx, userID)
	}
}

func handleGateTurn(
	ctx context.Context,
	w *sse.Writer,
	exec *session.Executor,
	st store.Store,
	sessionID, userID string,
	in session.Input,
	gate safety.Result,
	userMessage string,
) {
	outcome, err := exec.RunTurn(ctx, in)
	if err != nil {
		writeProviderError(w, exec, err)
		return
	}

	if outcome.Crisis != nil {
		safety.Audit(sessionID, gate, userMessage)
		_ = session.PersistTurn(ctx, st, sessionID, userID, gate, userMessage, outcome.Crisis.Body)
		writeTemplateEvent(w, "crisis", outcome.Crisis.TemplateID, outcome.Crisis.Body)
		return
	}
	if outcome.Medical != nil {
		safety.Audit(sessionID, gate, userMessage)
		_ = session.PersistTurn(ctx, st, sessionID, userID, gate, userMessage, outcome.Medical.Body)
		writeTemplateEvent(w, "medical", outcome.Medical.TemplateID, outcome.Medical.Body)
		return
	}
	if outcome.Block != nil {
		safety.Audit(sessionID, gate, userMessage)
		_ = session.PersistTurn(ctx, st, sessionID, userID, gate, userMessage, "")
		writeGateError(w, outcome.Block.Code, outcome.Block.Message)
		return
	}
	writeProviderError(w, exec, fmt.Errorf("unexpected pass outcome on gate branch"))
}

func writeDone(w *sse.Writer, sessionID string, st store.Store, ctx context.Context, userID string) {
	count := 0
	if sess, err := st.GetSession(ctx, sessionID); err == nil && sess != nil {
		count = sess.MessageCount()
	}
	done, err := json.Marshal(donePayload{SessionID: sessionID, MessageCount: count})
	if err != nil {
		return
	}
	if err := w.WriteEvent("", "done", done); err != nil {
		fmt.Printf("sse done write: %v\n", err)
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

func writeProviderError(w *sse.Writer, exec *session.Executor, err error) {
	msg := chatmodel.PublicMessage(err, exec.LLMConfig.ARKAPIKey)
	if errors.Is(err, chatmodel.ErrFailoverNotConfigured) {
		writeGateError(w, "failover_unavailable", msg)
		return
	}
	writeGateError(w, "provider_failed", msg)
}

func writeStreamError(w *sse.Writer, err error) {
	writeGateError(w, "stream_failed", err.Error())
}

func writeSessionHTTPError(c *app.RequestContext, err error) {
	switch {
	case errors.Is(err, store.ErrSessionNotFound):
		c.JSON(consts.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, store.ErrSessionForbidden):
		c.JSON(consts.StatusForbidden, map[string]string{"error": err.Error()})
	case errors.Is(err, store.ErrSessionEnded):
		c.JSON(consts.StatusConflict, map[string]string{"error": err.Error(), "code": "session_ended"})
	default:
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
