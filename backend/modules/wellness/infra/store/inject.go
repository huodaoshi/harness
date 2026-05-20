package store

import (
	"strings"

	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
)

// BuildInjectBlock formats profile + latest summary for prompt injection.
// Returns empty string when both are absent (caller may use placeholder).
func BuildInjectBlock(profile *domain.RelationshipProfile, summary *domain.SessionSummary) string {
	var b strings.Builder
	wrote := false

	if profile != nil && !profileIsEmpty(profile) {
		b.WriteString("[????]\n")
		if profile.Self.Note != "" {
			b.WriteString("???: ")
			b.WriteString(profile.Self.Note)
			b.WriteByte('\n')
		}
		for _, p := range profile.People {
			b.WriteString(p.Label)
			if p.Relation != "" {
				b.WriteString("?")
				b.WriteString(p.Relation)
				b.WriteString("?")
			}
			if p.Note != "" {
				b.WriteString(": ")
				b.WriteString(p.Note)
			}
			b.WriteByte('\n')
		}
		if profile.CurrentIssue != "" {
			b.WriteString("????: ")
			b.WriteString(profile.CurrentIssue)
			b.WriteByte('\n')
		}
		wrote = true
	}

	if summary != nil && len(summary.Summary3) > 0 {
		if wrote {
			b.WriteByte('\n')
		}
		b.WriteString("[??????]\n")
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

func profileIsEmpty(p *domain.RelationshipProfile) bool {
	if p.Self.Note != "" || p.CurrentIssue != "" || len(p.People) > 0 {
		return false
	}
	return true
}
