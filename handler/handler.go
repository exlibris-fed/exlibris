package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/model"

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
	var response []model.Book
	for _, book := range books {
		b := model.Book{
			Key:       book.Key,
			Title:     book.Title,
			Authors:   []model.Author{},
			Published: book.FirstPublishYear,
			//ISBN:      book.ISBN, // need to dedupe
			Subjects: book.Subject,
		}
		for _, a := range book.AuthorName {
			b.Authors = append(b.Authors, model.Author{Name: a})
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
