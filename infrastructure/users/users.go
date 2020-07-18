// Package users contains the repository for User objects.
package users

import (
	"errors"
	"strings"

	"github.com/exlibris-fed/exlibris/infrastructure/registrationkeys"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is the error returned when a user does not exist.
	ErrNotFound = errors.New("user could not be found")

	// ErrStorage is a generic database error.
	ErrStorage = errors.New("error with storage")

	// ErrDuplicate is the error returned when trying to create a user that already exists.
	ErrDuplicate = errors.New("user already exists")
)

// New returns a User Repository.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db:      db,
		keyRepo: registrationkeys.New(db),
	}
}

// A Repository is the repository pattern for Users.
type Repository struct {
	db      *gorm.DB
	keyRepo *registrationkeys.Repository
}

// GetByUsername returns a User object given a username. It does not fill in any related objects via `Preload`.
func (r *Repository) GetByUsername(name string) (*model.User, error) {
	var user model.User
	result := r.db.Where("username = ?", name).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return &user, nil
}

// GetByUsernameWithFollowers returns a User object given a username. It includes their list of followers.
func (r *Repository) GetByUsernameWithFollowers(name string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Followers").Where("username = ?", name).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return &user, nil
}

// Create persists a new user.
func (r *Repository) Create(user *model.User, key *model.RegistrationKey) (*model.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := tx.Create(key).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, ErrDuplicate
		}
	}
	return user, err
}

// Save will persist updates to a user.
func (r *Repository) Save(user *model.User) (*model.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := r.db.Save(user)
		if result.Error != nil {
			return ErrStorage
		}
		user = result.Value.(*model.User)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Activate will set a user's verified status to true, removing the registration key.
func (r *Repository) Activate(id uuid.UUID) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		key, err := r.keyRepo.Get(id)
		if err != nil {
			return err
		}
		key.User.Verified = true
		if _, err := r.Save(&key.User); err != nil {
			return err
		}
		if err := r.keyRepo.Delete(key); err != nil {
			return err
		}
		return nil
	})
}
