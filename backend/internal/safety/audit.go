package safety

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"
)

// Audit logs gate outcomes without storing raw user message text.
// Fields align with optional audit_events: gate_result, session_id, timestamp.
func Audit(sessionID string, r Result, message string) {
	hash := ""
	if message != "" {
		hash = HashMessage(message)
	}
	slog.Info("audit_event",
		"session_id", sessionID,
		"gate_result", string(r.Level),
		"template_id", r.TemplateID,
		"message_hash", hash,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// HashMessage returns a short hash for optional diagnostics (not reversible).
func HashMessage(message string) string {
	sum := sha256.Sum256([]byte(message))
	return hex.EncodeToString(sum[:8])
}
