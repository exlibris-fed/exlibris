package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
)

// AccountURIScheme is the string "acct" as specified in RFC 7565
const AccountURIScheme = "acct"

const (
	// RelWebfingerProfilePage is the link to Webfinger's profile page
	RelWebfingerProfilePage = "http://webfinger.net/rel/profile-page"

	// RelSelf is the "self" rel link
	RelSelf = "self"
)

// HandleWebfinger handles Webfinger requests (https://tools.ietf.org/html/rfc7033)
// to look up a user by their profile id. It is required for Mastodon interoperability.
//
// TODO: this provides basic functionality; I haven't actually looked at the RFC.
func (h *Handler) HandleWebfinger(w http.ResponseWriter, r *http.Request) {
	resource := r.FormValue("resource")
	if resource == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	split := strings.Split(resource, ":")
	if len(split) != 2 || split[0] != AccountURIScheme {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := strings.Split(split[1], "@")
	if len(username) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if username[1] != h.cfg.Domain {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var user model.User
	if err := h.db.Where("username = ?", username[0]).
		First(&user).Error; err != nil {
		// TODO may not be a 404
		w.WriteHeader(http.StatusNotFound)
		return
	}

	atProfile := fmt.Sprintf("%s://%s/@%s", h.cfg.Scheme, h.cfg.Domain, user.Username)
	userProfile := fmt.Sprintf("%s://%s/user/%s", h.cfg.Scheme, h.cfg.Domain, user.Username)
	response := dto.Webfinger{
		Subject: resource,
		Aliases: []string{atProfile, userProfile},
		Links: []dto.WebfingerLink{
			dto.WebfingerLink{
				Rel:  RelWebfingerProfilePage,
				Type: "text/html",
				Href: atProfile,
			},
			dto.WebfingerLink{
				Rel:  RelSelf,
				Type: "application/activity+json",
				Href: userProfile,
			},
		},
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshalling json for webfinger user %s: %s", user.Username, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
