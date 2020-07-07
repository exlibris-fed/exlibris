package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
	//"github.com/exlibris-fed/openlibrary-go"
	//"github.com/gorilla/mux"
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
	/*
		vars := mux.Vars(r)
		id := vars["book"]
		c := r.Context()
		username, ok := c.Value(model.ContextKeyAuthenticatedUsername).(string)
		if !ok {
			// the middleware should have required this already
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		book := model.Book{}
		h.db.Where("id = ?", fmt.Sprintf("/works/%s", id)).First(&book)
		if book.ID == "" {
			// fetch book from API
			work, err := openlibrary.GetWorkByID(id)
			if err != nil {
				log.Println("Could not fetch work", id, "got error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("%+v\n", work)
			return

				book := model.NewBook(work.Key, work.Title, 0, "")
				result := h.db.Create(book)
				if result.Error != nil {
					log.Println("Could not insert book into DB")
					w.WriteHeader(http.StatusConflict)
					return
				}
		}

		read := model.Read{
			UserID: user.ID,
			User:   &user,
			BookID: id,
			Book:   &book,
		}
		h.db.Create(&read)

		log.Printf("%+v", read)

		//go h.ap.Federate(c, user, read)

		w.WriteHeader(http.StatusCreated)
	*/
}
