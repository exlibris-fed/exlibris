package reads

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound   = errors.New("reads could not be found for user")
	ErrNotCreated = errors.New("read could not be saved")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetByUser(user *model.User) ([]*model.Read, error) {
	reads := []*model.Read{}
	result := r.db.Preload("Book").
		Preload("Book.Authors").
		Where("user_id = ?", user.ID).
		Order("created_at desc").
		Find(&reads)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return reads, nil
}

func (r *Repository) Create(read *model.Read) (*model.Read, error) {
	result := r.db.Create(read)
	if result.Error != nil {
		return nil, ErrNotCreated
	}
	return result.Value.(*model.Read), nil
}

// Get retrieves a read by its id (which is a uri to the activity).
func (r *Repository) Get(id string) (result *model.Read, err error) {
	result = new(model.Read)
	err = r.db.Where("id = ?", id).
		First(result).
		Error
	return
}
