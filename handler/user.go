package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/key"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/gorilla/mux"
)

// HandleActivityPubProfile returns a user's profile when requested with the AP content type.
func (h *Handler) HandleActivityPubProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var user model.User
	if err := h.db.Where("username = ?", username).
		First(&user).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("error retrieving user %s: %s", username, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	profile := fmt.Sprintf("%s/user/%s", h.cfg.Domain, user.Username)
	response := dto.NewActivityPubUser()
	response.ID = profile
	response.Following = profile + "/following"
	response.Followers = profile + "/followers"
	response.Inbox = profile + "/inbox"
	response.Outbox = profile + "/outbox"
	response.Username = user.Username
	response.Name = user.DisplayName
	response.URL = fmt.Sprintf("%s/@%s", h.cfg.Domain, user.Username)
	response.Endpoints["sharedInbox"] = fmt.Sprintf("%s/inbox", h.cfg.Domain)

	if publicKey, err := marshalPublicKey(user.PrivateKey); err == nil {
		response.PublicKey = dto.PublicKey{
			ID:    profile + "#main-key",
			Owner: profile,
			PEM:   publicKey,
		}
	} else {
		log.Printf("unable to marshal public key for user %s: %s", user.Username, err.Error())
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshalling json for user %s profile: %s", user.Username, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/activity+json")
	w.Write(b)
}

func marshalPublicKey(b []byte) (string, error) {
	privateKey, err := key.DeserializeRSAPrivateKey(b)
	if err != nil {
		return "", err
	}
	return key.MarshalPublicKeyFromPrivateKey(privateKey)
}
