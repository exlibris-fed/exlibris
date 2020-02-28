package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/exlibris-fed/openlibrary-go"
	"github.com/gorilla/mux"
)

// GetReads returns a list of books you've read
func (h *Handler) GetReads(w http.ResponseWriter, r *http.Request) {
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
			ID:        book.Key,
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

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["book"]
	c := h.contextFromRequest(r)
	userctx := c.Value(model.ContextKeyAuthenticatedUser)
	if userctx == nil {
		log.Println("User is nil")
	}
	user := userctx.(model.User)

	log.Println(user.Username)
	log.Println(id)
	book := model.Book{}
	h.db.Where("key = ?", fmt.Sprintf("works/%s",id)).First(&book)
	log.Println(book.Key)
	if book.Key == "" {
		// fetch book from API
		work, err := openlibrary.GetWorkByID(id)
		if err != nil {
			log.Println("Could not fetch work", id, "got error", err)
			return
		}
		book := model.NewBook(work.Key, work.Title, 0, "")
		result := h.db.Create(book)
		if result.Error != nil {
			spew.Dump(result.Error)
		}
	} else {
		log.Println("Book found")
	}
	spew.Dump(book)
}