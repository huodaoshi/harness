package safety

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
)

// Audit logs gate outcomes without storing raw user message text.
func Audit(sessionID string, r Result) {
	hash := ""
	if r.IsCrisis() {
		// Placeholder: no message body available here; never log raw text.
		hash = "redacted"
	}
	slog.Info("safety_gate",
		"session_id", sessionID,
		"gate_result", string(r.Level),
		"template_id", r.TemplateID,
		"message_hash", hash,
	)
}

// HashMessage returns a short hash for optional diagnostics (not reversible).
func HashMessage(message string) string {
	sum := sha256.Sum256([]byte(message))
	return hex.EncodeToString(sum[:8])
}
