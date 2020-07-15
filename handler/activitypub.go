package handler

import (
	"log"
	"net/http"
)

func (h *Handler) HandleActivityPubAction(w http.ResponseWriter, r *http.Request) {
	log.Println("handling ap action")

	if isActivityPubRequest, err := h.streamHandler(r.Context(), w, r); err != nil {
		// Do something with `err`
		log.Println("error handling ActivityStreams request:", err.Error())
		return
	} else if isActivityPubRequest {
		// Go-Fed handled the ActivityPub GET request for this particular IRI
		log.Println("go fed took care of it")
		return
	}
	// Here we return an error, but you may just as well decide
	// to render a webpage instead. But be sure you've already
	// applied the appropriate authorizations.
	log.Println("wtf...?")
	http.Error(w, "Non-ActivityPub request", http.StatusBadRequest)
	return
}
