package handler

import (
	"github.com/exlibris-fed/exlibris/activitypub"
	"github.com/exlibris-fed/exlibris/config"
	"github.com/exlibris-fed/exlibris/infrastructure/authors"
	"github.com/exlibris-fed/exlibris/infrastructure/books"
	"github.com/exlibris-fed/exlibris/infrastructure/reads"
	"github.com/exlibris-fed/exlibris/infrastructure/registrationkeys"
	"github.com/exlibris-fed/exlibris/infrastructure/reviews"
	"github.com/exlibris-fed/exlibris/infrastructure/users"
	"github.com/exlibris-fed/exlibris/service"

	"github.com/go-fed/activity/pub"
	"github.com/jinzhu/gorm"
)

// A Handler accepts http requests.
type Handler struct {
	ap                   *activitypub.ActivityPub
	actor                pub.FederatingActor
	streamHandler        pub.HandlerFunc
	cfg                  *config.Config
	bookService          *service.Book
	authorService        *service.Author
	editionsService      *service.Editions
	booksRepo            *books.Repository
	reviewsRepo          *reviews.Repository
	authorsRepo          *authors.Repository
	usersRepo            *users.Repository
	readsRepo            *reads.Repository
	registrationKeysRepo *registrationkeys.Repository
}

// New creates a new Handler to be used in processing http requests.
func New(db *gorm.DB, cfg *config.Config) *Handler {
	ap := activitypub.New(db)
	return &Handler{
		cfg:                  cfg,
		ap:                   activitypub.New(db),
		actor:                ap.NewFederatingActor(),
		streamHandler:        ap.NewStreamsHandler(),
		bookService:          service.NewBook(db),
		authorService:        service.NewAuthor(db),
		editionsService:      service.NewEditions(db),
		booksRepo:            books.New(db),
		reviewsRepo:          reviews.New(db),
		authorsRepo:          authors.New(db),
		usersRepo:            users.New(db),
		readsRepo:            reads.New(db),
		registrationKeysRepo: registrationkeys.New(db),
	}
}
