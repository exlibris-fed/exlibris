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
	c := h.contextFromRequest(r)
	userctx := c.Value(model.ContextKeyAuthenticatedUser)
	if userctx == nil {
		log.Println("User is nil")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := userctx.(model.User)
	reads := []model.Read{}
	response := []dto.Book{}
	h.db.Where("fk_user = ?", user.ID).Find(&reads)
	for _, read := range reads {
		log.Println(read.FKBook)
		book := model.Book{}
		h.db.Where("key = ?", fmt.Sprintf("/works/%s",read.FKBook)).First(&book)
		spew.Dump(book)
		response = append(response, dto.Book{ID: book.Key,Title: book.Title})
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	id := vars["book"]
	c := h.contextFromRequest(r)
	userctx := c.Value(model.ContextKeyAuthenticatedUser)
	if userctx == nil {
		log.Println("User is nil")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := userctx.(model.User)
	book := model.Book{}
	h.db.Where("key = ?", fmt.Sprintf("/works/%s",id)).First(&book)
	if book.Key == "" {
		// fetch book from API
		work, err := openlibrary.GetWorkByID(id)
		if err != nil {
			log.Println("Could not fetch work", id, "got error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		book := model.NewBook(work.Key, work.Title, 0, "")
		result := h.db.Create(book)
		if result.Error != nil {
			log.Println("Could not insert book into DB")
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	read := model.Read{
		FKUser: user.ID,
		FKBook: id,
	}
	h.db.Create(&read)
	w.WriteHeader(http.StatusCreated)
}