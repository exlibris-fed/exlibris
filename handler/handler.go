package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
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

// Register will create a user on the server
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var request dto.Register
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.DisplayName == "" || request.Email == "" || request.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := model.NewUser(request.Username, request.Password, request.Email, request.DisplayName)
	if err != nil {
		log.Println("error creating user object: " + err.Error())
	}
	result := h.db.Create(user)

	if result.Error != nil {
		// I'd like this to be a constant in a db package somewhere
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Authenticate will validate a user's password and, if correct, return a JWT
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var request dto.AuthenticationRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if request.Username == "" || request.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user model.User
	h.db.First(&user, "username = ?", request.Username)
	if user.PrivateKey == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !user.IsPassword(request.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := user.GenerateJWT()
	if err != nil {
		log.Printf("error generating jwt for user %s: %s", user.Username, err)
		// still return 401 because auth failed. is this correct?
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := dto.AuthenticationResponse{
		JWT: jwt,
	}
	b, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling jwt json: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write(b)
}

func (h *Handler) contextFromRequest(r *http.Request) context.Context {
	vars := mux.Vars(r)

	token := jwtFromAuth(r.Header.Get("Authorization"))
	c := context.WithValue(context.Background(), model.ContextKeyJWT, token)

	var user model.User
	h.db.First(&user, "username = ?", kidFromJWT(r.Header.Get("Authorization")))
	if len(user.PrivateKey) > 0 && user.ValidateJWT(token) {
		c = context.WithValue(c, model.ContextKeyAuthenticatedUser, user)
	}

	username, ok := vars["username"]
	if !ok {
		return c
	}
	c = context.WithValue(c, model.ContextKeyRequestedUser, username)
	return c
}

func (h *Handler) FederationTest(w http.ResponseWriter, r *http.Request) {
}

// HandleInbox is the http handler for an ActivityPub user's inbox.
func (h *Handler) HandleInbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin inbox")
	actor := h.ap.NewFederatingActor()

	c := h.contextFromRequest(r)

	if handled, err := actor.PostInbox(c, w, r); err != nil {
		log.Printf("error handling PostInbox: %s", err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		return
	} else if handled {
		log.Println("handled PostInbox")
		return
	} else if handled, err = actor.GetInbox(c, w, r); err != nil {
		log.Printf("error handling GetInbox: %s", err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		// Write to w
		return
	} else if handled {
		log.Println("handled GetInbox")
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

	c := h.contextFromRequest(r)

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
