package authors

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound = errors.New("author could not be found")
	// ErrNotCreated is returned when a record cannot be created
	ErrNotCreated = errors.New("author could not be created")
)

// New creates a new Repository instance for authors
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Repository is used for querying and creating authors
type Repository struct {
	db *gorm.DB
}

// GetByID returns an author from the database given an ID
// Will also return their books preloaded
func (r *Repository) GetByID(id string) (*model.Author, error) {
	var author model.Author
	result := r.db.Preload("Books").
		Where("open_library_id = ?", id).
		First(&author)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return &author, nil
}

// Create will create an author in the database from a model
func (r *Repository) Create(author *model.Author) (*model.Author, error) {
	result := r.db.Create(author)
	if result.Error != nil {
		return nil, ErrNotCreated
	}
	return result.Value.(*model.Author), nil
}
