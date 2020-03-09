package handler

import (
	"net/http"

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


func (h *Handler) FederationTest(w http.ResponseWriter, r *http.Request) {
}


