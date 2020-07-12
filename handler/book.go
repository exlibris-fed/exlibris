package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/gorilla/mux"
)

var sizeMapping = map[string]string{
	"S": "small",
	"M": "medium",
	"L": "large",
}

// GetBook returns a book
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["book"]

	book, err := h.bookService.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := dto.Book{
		ID:          book.OpenLibraryID,
		Title:       book.Title,
		Description: book.Description,
		Published:   time.Unix(int64(book.Published), 0),
	}
	response.Covers = make(map[string]string)
	for _, cover := range book.Covers {
		response.Covers[sizeMapping[cover.Type]] = cover.URL
	}
	for _, author := range book.Authors {
		response.Authors = append(response.Authors, author.Name)
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling json: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
