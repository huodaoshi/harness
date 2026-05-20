package store

import (
	"time"

	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
)

// EmptyProfile is the GET response when no profile exists yet (200, not 404).
func EmptyProfile(userID string) domain.RelationshipProfile {
	return domain.RelationshipProfile{
		UserID:       userID,
		Self:         domain.ProfileSelf{Note: ""},
		People:       []domain.ProfilePerson{},
		CurrentIssue: "",
	}
}

// NormalizeProfile ensures people slice is non-nil for JSON clients.
func NormalizeProfile(p *domain.RelationshipProfile) domain.RelationshipProfile {
	if p == nil {
		return EmptyProfile("")
	}
	out := *p
	if out.People == nil {
		out.People = []domain.ProfilePerson{}
	}
	if out.UpdatedAt.IsZero() {
		out.UpdatedAt = time.Time{}
	}
	return out
}
