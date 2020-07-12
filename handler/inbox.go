package handler

import (
	"log"
	"net/http"
)

// HandleInbox is the http handler for an ActivityPub user's inbox.
func (h *Handler) HandleInbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin inbox")

	c := r.Context()

	if handled, err := h.actor.PostInbox(c, w, r); err != nil {
		log.Printf("error handling PostInbox: %s", err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		return
	} else if handled {
		log.Println("handled PostInbox")
		return
	} else if handled, err = h.actor.GetInbox(c, w, r); err != nil {
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
