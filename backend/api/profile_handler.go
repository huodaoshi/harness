package api

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/huodaoshi/harness/backend/modules/wellness/domain"
	wellnessstore "github.com/huodaoshi/harness/backend/modules/wellness/infra/store"
)

// profilePutBody is the PUT /v1/profile JSON body (user_id from auth, not body).
type profilePutBody struct {
	Self         domain.ProfileSelf     `json:"self"`
	People       []domain.ProfilePerson `json:"people"`
	CurrentIssue string                 `json:"current_issue"`
}

// NewGetProfileHandler returns GET /v1/profile for the authenticated user_id.
func NewGetProfileHandler(st domain.Store) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := requireUserID(c)
		if !ok {
			return
		}

		p, err := st.GetProfile(ctx, userID)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if p == nil {
			c.JSON(consts.StatusOK, wellnessstore.EmptyProfile(userID))
			return
		}
		c.JSON(consts.StatusOK, p)
	}
}

// NewPutProfileHandler returns PUT /v1/profile (full upsert).
func NewPutProfileHandler(st domain.Store) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := requireUserID(c)
		if !ok {
			return
		}

		var body profilePutBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]string{"error": "invalid json body"})
			return
		}
		if body.People == nil {
			body.People = []domain.ProfilePerson{}
		}

		p := domain.RelationshipProfile{
			UserID:       userID,
			Self:         body.Self,
			People:       body.People,
			CurrentIssue: body.CurrentIssue,
			UpdatedAt:    time.Now().UTC(),
		}
		if err := st.UpsertProfile(ctx, p); err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		saved, err := st.GetProfile(ctx, userID)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(consts.StatusOK, saved)
	}
}
