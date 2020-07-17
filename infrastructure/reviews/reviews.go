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
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound = errors.New("reviews could not be found for book")
	// ErrNotCreated is returned when a record cannot be created
	ErrNotCreated = errors.New("review for book could not be created")
	// ErrStorage is returned when an unknown storage issue occurs
	ErrStorage = errors.New("could not retrieve data")
)

// New creates a new Repository instance for reviews
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Repository is used for querying and creating reviews
type Repository struct {
	db *gorm.DB
}

// GetReviews returns the reviews for a given book
// Preloads the User object
func (r *Repository) GetReviews(book *model.Book) ([]model.Review, error) {
	var reviews []model.Review
	if err := r.db.Preload("User").
		Where("book_id = ?", book.OpenLibraryID).
		Find(&reviews).Error; err != nil {
		return nil, ErrNotFound
	}
	return reviews, nil
}

// CreateReview will create a new review for a given book
func (r *Repository) CreateReview(user *model.User, book *model.Book, text string, rating int) (*model.Review, error) {
	// @TODO: rating
	book, err := books.New(r.db).GetByID(book.OpenLibraryID)
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
