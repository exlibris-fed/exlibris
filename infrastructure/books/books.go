package books

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound = errors.New("book could not be found")
	// ErrNotCreated is returned when a record cannot be created
	ErrNotCreated = errors.New("book could not be created")
	// ErrStorage is returned when an unknown storage issue occurs
	ErrStorage = errors.New("error with storage")
)

// New creates a new Repository instance for books
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Repository is used for querying and creating books
type Repository struct {
	db *gorm.DB
}

// GetByID returns a book from the database given an ID
// Will also return its authors and covers
func (r *Repository) GetByID(id string) (*model.Book, error) {
	var book model.Book
	result := r.db.Preload("Covers").
		Preload("Authors").
		Where("open_library_id = ?", id).
		First(&book)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}

	return &book, nil
}

func (r *Repository) Create(book *model.Book) (*model.Book, error) {
	result := r.db.Create(book)
	if result.Error != nil {
		return nil, ErrNotCreated
	}
	return result.Value.(*model.Book), nil
}
