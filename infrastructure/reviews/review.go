package reviews

import (
	"errors"
	"fmt"
	"log"

	"github.com/exlibris-fed/exlibris/infrastructure/books"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound   = errors.New("reviews could not be found for book")
	ErrNotCreated = errors.New("review for book could not be created")
	ErrStorage    = errors.New("could not retrieve data")
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
	book, err := books.New(r.db).GetByID("/works/" + id)
	if err != nil {
		if errors.Is(err, books.ErrNotFound) {
			return nil, fmt.Errorf("error book does not exist: %w", ErrNotFound)
		}
		log.Println("could not retrieve book", err)
		return nil, fmt.Errorf("error finding book: %w", ErrStorage)
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
