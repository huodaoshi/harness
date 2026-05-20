package store

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
)

// GenerateSummary3 builds a 3-line post-session summary (MVP: no LLM).
func GenerateSummary3(sess *domain.StoredSession) []string {
	if sess == nil {
		return []string{
			"????????",
			"???????",
			"?????????",
		}
	}

	var userSnippets []string
	for _, m := range sess.Messages {
		if m.Role == "user" && strings.TrimSpace(m.Content) != "" {
			userSnippets = append(userSnippets, strings.TrimSpace(m.Content))
		}
	}

	modeLabel := "????"
	if sess.Mode == "distress" {
		modeLabel = "????"
	}

	if len(userSnippets) == 0 {
		return []string{
			fmt.Sprintf("????%s??????", modeLabel),
			"???????????",
			"??????????????",
		}
	}

	last := userSnippets[len(userSnippets)-1]
	return []string{
		fmt.Sprintf("??????%s", truncateRunes(last, 36)),
		fmt.Sprintf("?????%s?", modeLabel),
		"????????????????????????",
	}
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max]) + "?"
}
