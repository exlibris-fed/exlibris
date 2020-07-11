package reviews

import (
	"errors"
	"fmt"

	"github.com/exlibris-fed/exlibris/infrastructure/books"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound   = errors.New("reviews could not be found for book")
	ErrNotCreated = errors.New("review for book could not be created")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetReviews(id string) ([]*model.Review, error) {
	var reviews []*model.Review
	if err := r.db.Preload("User").Where(&reviews, "book_id = ?", id).Find(&reviews).Error; err != nil {
		return nil, ErrNotFound
	}
	return reviews, nil
}

func (r *Repository) CreateReview(user *model.User, id string, text string, rating int) (*model.Review, error) {
	// @TODO: rating
	book, err := books.New(r.db).GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding book: %w", ErrNotCreated)
	}

	review := &model.Review{
		Base: model.Base{
			ID: uuid.New(),
		},
		Book:   *book,
		BookID: book.OpenLibraryID,
		Text:   text,
		User:   *user,
		UserID: user.ID,
	}

	if result := r.db.Create(&review); result.Error != nil {
		return nil, fmt.Errorf("error saving book: %w", ErrNotCreated)
	}

	return review, nil
}
