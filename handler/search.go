package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/openlibrary-go"
)

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
			ID: book.Key,
			Title:     book.Title,
			Authors:   []string{},
			Published: book.FirstPublishYear,
			//ISBN:      book.ISBN, // need to dedupe
			Subjects: book.Subject,
			Covers: map[string]string{
				"small":  book.CoverURL(openlibrary.SizeSmall),
				"medium": book.CoverURL(openlibrary.SizeMedium),
				"large":  book.CoverURL(openlibrary.SizeLarge),
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
