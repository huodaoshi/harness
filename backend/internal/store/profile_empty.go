package store

import "time"

// EmptyProfile is the GET response when no profile exists yet (200, not 404).
func EmptyProfile(userID string) RelationshipProfile {
	return RelationshipProfile{
		UserID:       userID,
		Self:         ProfileSelf{Note: ""},
		People:       []ProfilePerson{},
		CurrentIssue: "",
	}
}

// NormalizeProfile ensures people slice is non-nil for JSON clients.
func NormalizeProfile(p *RelationshipProfile) RelationshipProfile {
	if p == nil {
		return EmptyProfile("")
	}
	out := *p
	if out.People == nil {
		out.People = []ProfilePerson{}
	}
	if out.UpdatedAt.IsZero() {
		out.UpdatedAt = time.Time{}
	}
	return out
}
