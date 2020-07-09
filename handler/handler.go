package handler

import (
	"github.com/exlibris-fed/exlibris/activitypub"
	"github.com/exlibris-fed/exlibris/config"
	"github.com/exlibris-fed/exlibris/service"

	"github.com/jinzhu/gorm"
)

// A Handler accepts http requests.
type Handler struct {
	db              *gorm.DB
	ap              *activitypub.ActivityPub
	cfg             *config.Config
	bookService     *service.Book
	authorService   *service.Author
	editionsService *service.Editions
}

// New creates a new Handler to be used in processing http requests.
func New(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:              db,
		cfg:             cfg,
		ap:              activitypub.New(db),
		bookService:     service.NewBook(db),
		authorService:   service.NewAuthor(db),
		editionsService: service.NewEditions(db),
	}
}
