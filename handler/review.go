package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/exlibris-fed/exlibris/dto"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReviewRequest struct {
	Review string `json:"review"`
}

func (h *Handler) Review(w http.ResponseWriter, r *http.Request) {
	log.Println("Review handler")
	c := r.Context()
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(model.User)

	if !ok {
		// error
		return
	}

	vars := mux.Vars(r)
	id := vars["book"]

	var review []model.Review
	var response []dto.Review
	var err error
	if r.Method == http.MethodGet {
		// get the review for the book
		review, err = h.getReview(id)
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

		review, err = h.createReview(&user, id, reviewData.Review)
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

	for _, review := range review {
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

func (h *Handler) getReview(id string) ([]model.Review, error) {
	var review []model.Review
	if err := h.db.Preload("User").Where(&review, "book_id = ?", id).Find(&review).Error; err != nil {
		return nil, fmt.Errorf("Could not find review: %w", err)
	}
	return review, nil

}

func (h *Handler) createReview(user *model.User, id string, text string) ([]model.Review, error) {
	// @TODO: rating
	book := h.bookService.Get(id)

	review := model.Review{
		Base: model.Base{
			ID: uuid.New(),
		},
		Book:   *book,
		BookID: book.OpenLibraryID,
		Text:   text,
		User:   *user,
		UserID: user.ID,
	}

	spew.Dump(review)

	if result := h.db.Debug().Create(&review); result.Error != nil {
		return nil, result.Error
	}

	return []model.Review{review}, nil
}
