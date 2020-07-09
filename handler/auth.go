package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/mail"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	key, err := model.NewRegistrationKey(*user)
	if err != nil {
		log.Println("error creating registration key object: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := tx.Create(key).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		// I'd like this to be a constant in a db package somewhere
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	m := mail.New(h.cfg.SMTP.Host, h.cfg.SMTP.Port, h.cfg.SMTP.Username, h.cfg.SMTP.Password)
	if err := m.SendVerificationEmail(user.Email, fmt.Sprintf("%s/verify/%s", h.cfg.Domain, key.Key.String())); err != nil {
		// the user was created, so this isn't an error, but it's not great
		log.Printf("error sending registration email to user %s: %s", user.Username, err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// ResendVerificationKey will resend the key to a user if they are unverified.
func (h *Handler) ResendVerificationKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["user"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var key model.RegistrationKey
	if err := h.db.Table("registration_keys").
		Preload("User").
		Joins("inner join users on registration_keys.user_id = users.id").
		Where("username = ?", username).
		First(&key).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("error getting key: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	m := mail.New(h.cfg.SMTP.Host, h.cfg.SMTP.Port, h.cfg.SMTP.Username, h.cfg.SMTP.Password)
	if err := m.SendVerificationEmail(key.User.Email, fmt.Sprintf("%s/verify/%s", h.cfg.Domain, key.Key.String())); err != nil {
		// the user was created, so this isn't an error, but it's not great
		log.Printf("error sending registration email to user %s: %s", key.User.Username, err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
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

	if !user.Verified {
		w.WriteHeader(http.StatusForbidden)
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
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// VerifyKey validates a user's verification key and activates their account
func (h *Handler) VerifyKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringKey, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key, err := uuid.Parse(stringKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var registrationKey model.RegistrationKey
	if err := h.db.Preload("User").
		Where("key = ?", key).
		First(&registrationKey).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("error getting key: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		registrationKey.User.Verified = true
		if err := tx.Save(&registrationKey.User).Error; err != nil {
			return err
		}
		if err := tx.Delete(&registrationKey).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Printf("error activating user %s: %s", registrationKey.User.Username, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
