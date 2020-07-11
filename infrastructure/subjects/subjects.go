package subjects

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound = errors.New("Subject could not be found for book")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetByID(id string) ([]*model.Subject, error) {
	return nil, nil
}

func (r *Repository) Create(subjects []*model.Subject, book *model.Book) ([]*model.Subject, error) {
	return nil, nil
}
