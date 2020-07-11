package authors

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound   = errors.New("author could not be found")
	ErrNotCreated = errors.New("author could not be created")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

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

func (r *Repository) Create(author *model.Author) (*model.Author, error) {
	result := r.db.Create(author)
	if result.Error != nil {
		return nil, ErrNotCreated
	}
	return result.Value.(*model.Author), nil
}
