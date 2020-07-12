package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/dto"
	reviewsinfra "github.com/exlibris-fed/exlibris/infrastructure/reviews"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/gorilla/mux"
)

type ReviewRequest struct {
	Review string `json:"review"`
}

func (h *Handler) Review(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(*model.User)

	if !ok {
		log.Println("user not found")
		w.WriteHeader(http.StatusUnauthorized)
		// error
		return
	}

	vars := mux.Vars(r)
	id := vars["book"]

	var reviews []*model.Review
	var response []dto.Review
	var err error
	if r.Method == http.MethodGet {
		// get the review for the book
		reviews, err = h.reviewsRepo.GetReviews(id)
		if err != nil {
			log.Println(err)
		}
	} else if r.Method == http.MethodPost {
		// Create review of book

		decoder := json.NewDecoder(r.Body)
		var reviewData ReviewRequest
		err = decoder.Decode(&reviewData)
		if err != nil {
			// error with request
			log.Println("Could not read review")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var review *model.Review
		review, err = h.reviewsRepo.CreateReview(user, id, reviewData.Review, 0)
		if errors.Is(err, reviewsinfra.ErrNotFound) {
			// Trying to create a review about a book no one has viewed or read
			w.WriteHeader(http.StatusNotFound)
			return
		}
		reviews = []*model.Review{review}
	} else {
		log.Println("Bad request")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, review := range reviews {
		response = append(response, dto.Review{
			Author:    review.User.DisplayName,
			Text:      review.Text,
			Timestamp: review.CreatedAt,
		})
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling json: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)

}
