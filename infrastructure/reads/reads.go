package reads

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is returned when a record cannot be found.
	ErrNotFound = errors.New("reads could not be found for user")
	// ErrNotCreated is returned when a record cannot be created.
	ErrNotCreated = errors.New("read could not be saved")
)

// New creates a new Repository instance for reads.
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Repository is used for querying and creating reads.
type Repository struct {
	db *gorm.DB
}

// Get returns reads from the database given a user.
// Will also return the books and its authors.
func (r *Repository) Get(user *model.User) ([]*model.Read, error) {
	reads := []*model.Read{}
	result := r.db.Preload("Book").
		Preload("Book.Authors").
		Preload("Book.Covers").
		Where("user_id = ?", user.ID).
		Order("created_at desc").
		Find(&reads)
	if result.Error != nil {
		return nil, ErrNotFound
	}

	return reads, nil
}

// Create will persist the read to the database.
func (r *Repository) Create(read *model.Read) (*model.Read, error) {
	result := r.db.Create(read)
	if result.Error != nil {
		return nil, ErrNotCreated
	}
	return result.Value.(*model.Read), nil
}

// GetByID retrieves a read by its id (which is a uri to the activity).
func (r *Repository) GetByID(id string) (result *model.Read, err error) {
	result = new(model.Read)
	err = r.db.Where("id = ?", id).
		First(result).
		Error
	return
}
