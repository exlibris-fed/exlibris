package handler

import (
	"log"
	"net/http"

	"github.com/exlibris-fed/openlibrary-go"
	"github.com/jinzhu/gorm"
)

// A Handler accepts non-ActivityPub http requests. It may be better to move the activitypub handlers here eventually.
type Handler struct {
	db *gorm.DB
}

// New creates a new Handler to be used in processing http requests.
func New(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// SearchBooks will search the Library of Congress api for books. Currently only supports title search.
func (h *Handler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	books, err := openlibrary.TitleSearch(title)
	if err != nil {
		log.Println("error searching for titles: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, book := range books {
		log.Println(book.Title)
	}
}
