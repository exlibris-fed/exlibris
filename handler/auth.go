package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
)

// Register will create a user on the server
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		return
	}
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

