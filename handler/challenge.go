package handler

import (
	"net/http"
	"os"
)

// HandleChallenge will return the value in the environ variable ACME_CHALLENGE for let's encrypt
func (h *Handler) HandleChallenge(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(os.Getenv("ACME_CHALLENGE")))
	w.WriteHeader(http.StatusOK);
}