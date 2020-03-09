package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/gorilla/mux"
)

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
