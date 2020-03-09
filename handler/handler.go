package handler

import (
	"github.com/exlibris-fed/exlibris/activitypub"

	"github.com/jinzhu/gorm"
)

// A Handler accepts http requests.
type Handler struct {
	db *gorm.DB
	ap *activitypub.ActivityPub
}

// New creates a new Handler to be used in processing http requests.
func New(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
		ap: activitypub.New(db),
	}
}

/*
func (h *Handler) FederationTest(w http.ResponseWriter, r *http.Request) {
	c := h.contextFromRequest(r)
	userI := c.Value(model.ContextKeyAuthenticatedUser)
	if userI == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, ok := userI.(model.User)
	if !ok {
		log.Printf("userI is %T, not model.User", userI)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO
	readActivity := &model.Read{
		UserID: user.ID,
		FKBook: "QWERTY",
	}
	c = context.WithValue(c, model.ContextKeyRead, readActivity)

	actor := h.ap.NewFederatingActor()
	actor.Send(c, user.OutboxIRI(), read)
}
*/
