package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/exlibris-fed/openlibrary-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetReads returns a list of books a user has read
func (h *Handler) GetReads(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	username, ok := c.Value(model.ContextKeyAuthenticatedUsername).(string)
	if !ok {
		// the middleware should have required this already
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reads := []model.Read{}
	response := []dto.Book{}
	h.db.Preload("User").
		Preload("Book").
		Where("username = ?", username).
		Find(&reads)
	for _, read := range reads {
		log.Println(read.Book)
		response = append(response, dto.Book{ID: read.Book.OpenLibraryID, Title: read.Book.Title})
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

	book := &model.Book{}
	h.db.Where("open_library_id = ?", fmt.Sprintf("/works/%s", id)).First(&book)
	if book.OpenLibraryID == "" {
		// fetch book from API
		work, err := openlibrary.GetWorkByID(id)
		if err != nil {
			log.Println("Could not fetch work", id, "got error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("%+v\n", work)

		book = model.NewBook(work)
		result := h.db.Create(book)
		if result.Error != nil {
			log.Println("Could not insert book into DB:", result.Error)
			w.WriteHeader(http.StatusConflict)
			return
		}
		book = result.Value.(*model.Book)
	}

	read := model.Read{
		Base: model.Base{
			ID: uuid.New(),
		},
		UserID: user.ID,
		User:   user,
		BookID: book.ID,
		Book:   *book,
	}
	h.db.Debug().Create(&read)

	log.Printf("%+v", read)

	//go h.ap.Federate(c, user, read)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
