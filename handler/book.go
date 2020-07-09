package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/gorilla/mux"
)

// GetBook returns a book
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["book"]

	book := h.bookService.Get(id)
	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := dto.Book{ID: book.OpenLibraryID, Title: book.Title, Description: book.Description}
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
