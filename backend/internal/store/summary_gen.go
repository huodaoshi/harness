package store

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// GenerateSummary3 builds a 3-line post-session summary (MVP: no LLM).
func GenerateSummary3(sess *StoredSession) []string {
	if sess == nil {
		return []string{
			"本次对话已结束。",
			"感谢你的信任。",
			"随时可以再来聊聊。",
		}
	}

	var userSnippets []string
	for _, m := range sess.Messages {
		if m.Role == "user" && strings.TrimSpace(m.Content) != "" {
			userSnippets = append(userSnippets, strings.TrimSpace(m.Content))
		}
	}

	modeLabel := "日常陪伴"
	if sess.Mode == "distress" {
		modeLabel = "洪峰陪伴"
	}

	if len(userSnippets) == 0 {
		return []string{
			fmt.Sprintf("本次以「%s」模式结束。", modeLabel),
			"你还没有发送具体内容。",
			"下次可以从任何感受开始说起。",
		}
	}

	last := userSnippets[len(userSnippets)-1]
	return []string{
		fmt.Sprintf("你最近谈到：%s", truncateRunes(last, 36)),
		fmt.Sprintf("会话模式：%s。", modeLabel),
		"下次可以继续从这里说起，我会记得这场对话的摘要。",
	}
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max]) + "…"
}
