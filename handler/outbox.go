package handler

import (
	"log"
	"net/http"
)

func (h *Handler) HandleOutbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin outbox")
	actor := h.ap.NewFederatingActor()

	c := r.Context()

	// Populate c with request-specific information
	if handled, err := actor.PostOutbox(c, w, r); err != nil {
		// Write to w
		log.Println("error posting to outbox:", err.Error())
		return
	} else if handled {
		log.Println("handled post outbox")
		return
	} else if handled, err = actor.GetOutbox(c, w, r); err != nil {
		// Write to w
		log.Println("error getting outbox:", err.Error())
		return
	} else if handled {
		log.Println("handled get outbox")
		return
	}
	// else:
	//
	// Handle non-ActivityPub request, such as serving a webpage.
}
