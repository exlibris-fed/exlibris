package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetReads returns a list of books a user has read
func (h *Handler) GetReads(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(*model.User)
	if !ok {
		log.Println("No user")
		// the middleware should have required this already
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := []dto.Read{}

	reads, err := h.readsRepo.GetByUser(user)

	if err != nil {
		// Error searching
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, read := range reads {
		bookDTO := dto.Read{
			Book: dto.Book{
				ID:          read.Book.OpenLibraryID,
				Title:       read.Book.Title,
				Published:   time.Unix(int64(read.Book.Published), 0),
				Description: read.Book.Description,
			},
			Timestamp: read.CreatedAt,
		}
		for _, author := range read.Book.Authors {
			bookDTO.Authors = append(bookDTO.Authors, author.Name)
		}
		response = append(response, bookDTO)
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

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		log.Println("Bad method")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["book"]
	c := r.Context()
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(*model.User)
	if !ok {
		// the middleware should have required this already
		log.Println("Could not get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	book, err := h.bookService.Get(id)

	if err != nil {
		log.Println("could not fetch book for read", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bookID := fmt.Sprintf("%s://%s/user/%s/%s", h.cfg.Scheme, h.cfg.Domain, strings.ToLower(user.Username), uuid.New())
	read := model.Read{
		ID:     bookID,
		User:   *user,
		Book:   *book,
		BookID: book.OpenLibraryID,
	}

	// TODO: we're going to want to actually create as part of the AP flow. That's nearly ready but I'd like to discuss how much to grab here vs there (I think either is fine, because we can populate the data here and when it checks if we have the book/author/subjects/etc in activitypub/database's Create we don't fetch them)
	h.readsRepo.Create(&read)

	/*
		if _, err := h.actor.Send(c, user.OutboxIRI(), read.ToType()); err != nil {
			log.Printf("error sending to outbox for read %s: %s", read.ID, err.Error())
		}
	*/

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
