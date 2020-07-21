package registrationkeys

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is returned when a record cannot be found.
	ErrNotFound = errors.New("registration key could not be found")
	// ErrNotCreated is returned when a record cannot be created.
	ErrNotCreated = errors.New("registration key could not be saved")
	// ErrStorage is returned when an unknown storage issue occurs.
	ErrStorage = errors.New("error with storage")
)

// New creates a new Repository instance for registration keys.
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Repository is used for querying and creating registration keys.
type Repository struct {
	db *gorm.DB
}

// GetByUser returns a registration key given a user.
func (r *Repository) GetByUser(user *model.User) (*model.RegistrationKey, error) {
	var key model.RegistrationKey
	result := r.db.Model(user).Related(&key)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return &key, nil
}

// Get returns a registration key given its ID.
// Preloads the user object.
func (r *Repository) Get(id uuid.UUID) (*model.RegistrationKey, error) {
	var registrationKey model.RegistrationKey
	result := r.db.Preload("User").
		Where("key = ?", id).
		First(&registrationKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return &registrationKey, nil
}

// Create persists the given registration key.
func (r *Repository) Create(key *model.RegistrationKey) (*model.RegistrationKey, error) {
	result := r.db.Create(key)
	if result.Error != nil {
		return nil, ErrNotCreated
	}
	return result.Value.(*model.RegistrationKey), nil
}

// Delete will eliminate the registration key from the database.
func (r *Repository) Delete(key *model.RegistrationKey) error {
	return r.db.Delete(key).Error
}
