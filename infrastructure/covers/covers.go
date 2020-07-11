package covers

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound = errors.New("Cover could not be found for book")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetByID(id string) (*model.Book, error) {
	return nil, nil
}

func (r *Repository) Create(cover *model.Cover, book *model.Book) (*model.Cover, error) {
	return nil, nil
}
