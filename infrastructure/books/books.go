package books

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"

	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound   = errors.New("book could not be found")
	ErrNotCreated = errors.New("book could not be created")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetByID(id string) (*model.Book, error) {
	var book model.Book
	result := r.db.Preload("Covers").
		Preload("Authors").
		Where("open_library_id = ?", "/works/"+id).
		First(&book)
	if result.Error != nil {
		return nil, ErrNotFound
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
