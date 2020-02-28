package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/jinzhu/gorm"

	"github.com/exlibris-fed/openlibrary-go"
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
	var response []dto.Book
	for _, book := range books {
		b := dto.Book{
			Title:     book.Title,
			Authors:   []string{},
			Published: book.FirstPublishYear,
			//ISBN:      book.ISBN, // need to dedupe
			Subjects: book.Subject,
			Covers: map[string]string {
				"small":book.CoverURL(openlibrary.SizeSmall),
				"medium":book.CoverURL(openlibrary.SizeMedium),
				"large":book.CoverURL(openlibrary.SizeLarge),
			},
		}
		for _, a := range book.AuthorName {
			b.Authors = append(b.Authors, a)
		}
		response = append(response, b)
	}
	b, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling json: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
