package chatmodel

// PublicMessage returns a client-safe error string (no secrets).
func PublicMessage(err error, secrets ...string) string {
	if err == nil {
		return ""
	}
	return SanitizeError(err, secrets...).Error()
}
