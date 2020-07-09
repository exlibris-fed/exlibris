package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetReads returns a list of books a user has read
func (h *Handler) GetReads(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(model.User)
	if !ok {
		log.Println("No user")
		// the middleware should have required this already
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reads := []model.Read{}
	response := []dto.Read{}

	if err := h.db.Where("user_id = ?", user.ID).Find(&reads).Error; err != nil {
		// Error searching
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	for _, read := range reads {
		var book model.Book
		if err := h.db.Preload("Authors").Where(&model.Book{OpenLibraryID: read.BookID}).First(&book).Error; err != nil {
			log.Println(err)
			continue
		}
		bookDTO := dto.Read{Book: dto.Book{ID: book.OpenLibraryID, Title: book.Title, Published: time.Unix(int64(book.Published), 0), Description: book.Description}, Timestamp: read.CreatedAt}
		for _, author := range book.Authors {
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
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(model.User)
	if !ok {
		// the middleware should have required this already
		log.Println("Could not get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	book := h.bookService.Get(id)

	read := model.Read{
		Base: model.Base{
			ID: uuid.New(),
		},
		UserID: user.ID,
		User:   user,
		BookID: book.OpenLibraryID,
		Book:   *book,
	}
	h.db.Create(&read)

	//go h.ap.Federate(c, user, read)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
