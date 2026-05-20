package logging_test

import (
	"log/slog"
	"testing"

	"github.com/huodaoshi/harness/backend/infra/logging"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"", slog.LevelInfo},
		{"nope", slog.LevelInfo},
	}
	for _, tc := range tests {
		if got := logging.ParseLevel(tc.in); got != tc.want {
			t.Fatalf("ParseLevel(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
