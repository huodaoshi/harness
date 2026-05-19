package store

import (
	"strings"
)

// BuildInjectBlock formats profile + latest summary for prompt injection.
// Returns empty string when both are absent (caller may use placeholder).
func BuildInjectBlock(profile *RelationshipProfile, summary *SessionSummary) string {
	var b strings.Builder
	wrote := false

	if profile != nil && !profileIsEmpty(profile) {
		b.WriteString("[关系档案]\n")
		if profile.Self.Note != "" {
			b.WriteString("关于我: ")
			b.WriteString(profile.Self.Note)
			b.WriteByte('\n')
		}
		for _, p := range profile.People {
			b.WriteString(p.Label)
			if p.Relation != "" {
				b.WriteString("（")
				b.WriteString(p.Relation)
				b.WriteString("）")
			}
			if p.Note != "" {
				b.WriteString(": ")
				b.WriteString(p.Note)
			}
			b.WriteByte('\n')
		}
		if profile.CurrentIssue != "" {
			b.WriteString("当前议题: ")
			b.WriteString(profile.CurrentIssue)
			b.WriteByte('\n')
		}
		wrote = true
	}

	if summary != nil && len(summary.Summary3) > 0 {
		if wrote {
			b.WriteByte('\n')
		}
		b.WriteString("[上次会话摘要]\n")
		for i, line := range summary.Summary3 {
			b.WriteString(strings.TrimSpace(line))
			if i < len(summary.Summary3)-1 {
				b.WriteByte('\n')
			}
		}
		wrote = true
	}

	return strings.TrimSpace(b.String())
}

func profileIsEmpty(p *RelationshipProfile) bool {
	if p.Self.Note != "" || p.CurrentIssue != "" || len(p.People) > 0 {
		return false
	}
	return true
}
