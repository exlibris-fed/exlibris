package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/exlibris-fed/exlibris/activitypub"
	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/exlibris-fed/openlibrary-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// A Handler accepts http requests.
type Handler struct {
	db *gorm.DB
	ap *activitypub.ActivityPub
}

// New creates a new Handler to be used in processing http requests.
func New(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
		ap: activitypub.New(db),
	}
}

// HandleInbox is the http handler for an ActivityPub user's inbox.
func (h *Handler) HandleInbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin inbox")
	actor := h.ap.NewFederatingActor()

	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		// how did this happen? I almost want to make it a 500
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.User
	h.db.First(&user, "username = ?", username)
	c := context.WithValue(context.Background(), model.ContextKeyUser, user)
	c = context.WithValue(c, model.ContextKeyJWT, jwtFromAuth(r.Header.Get("Authorization")))

	if handled, err := actor.PostInbox(c, w, r); err != nil {
		log.Printf("error handling PostInbox for user %s: %s", username, err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		return
	} else if handled {
		log.Printf("handled PostInbox for user %s", username)
		return
	} else if handled, err = actor.GetInbox(c, w, r); err != nil {
		log.Printf("error handling GetInbox for user %s: %s", username, err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		// Write to w
		return
	} else if handled {
		log.Printf("handled GetInbox for user %s", username)
		return
	}
	log.Println("else...?")
	// else:
	//
	// Handle non-ActivityPub request, such as serving a webpage.
}

func (h *Handler) HandleOutbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin outbox")
	actor := h.ap.NewFederatingActor()

	// TODO
	c := context.Background()
	// Populate c with request-specific information
	if handled, err := actor.PostOutbox(c, w, r); err != nil {
		// Write to w
		return
	} else if handled {
		return
	} else if handled, err = actor.GetOutbox(c, w, r); err != nil {
		// Write to w
		return
	} else if handled {
		return
	}
	// else:
	//
	// Handle non-ActivityPub request, such as serving a webpage.
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
	var response []dto.Book
	for _, book := range books {
		b := dto.Book{
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

func jwtFromAuth(auth string) string {
	pieces := strings.Split(auth, " ")
	if len(pieces) != 2 {
		return ""
	}
	if strings.ToLower(pieces[0]) != "bearer" {
		return ""
	}
	return pieces[1]
}

func kidFromJWT(token string) string {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return ""
	}

	bytes, err := jwt.DecodeSegment(segments[1])
	if err != nil {
		log.Println("error decoding jwt segment: " + err.Error())
		return ""
	}
	var kid struct {
		ID string `json:"kid"`
	}
	err = json.Unmarshal(bytes, &kid)
	if err != nil {
		log.Println("error marshalling kid: " + err.Error())
		return ""
	}
	return kid.ID
}
