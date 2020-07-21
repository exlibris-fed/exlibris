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
	// ErrNotFound is returned when a record cannot be found.
	ErrNotFound = errors.New("user could not be found")
	// ErrNotCreated is returned when a record cannot be created.
	ErrNotCreated = errors.New("user could not be created")
	// ErrStorage is returned when an unknown storage issue occurs.
	ErrStorage = errors.New("error with storage")
	// ErrDuplicate occurs when the user already exists.
	ErrDuplicate = errors.New("user already exists")
)

// New creates a new Repository instance for users.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db:      db,
		keyRepo: registrationkeys.New(db),
	}
}

// Repository is used for querying and creating users.
type Repository struct {
	db      *gorm.DB
	keyRepo *registrationkeys.Repository
}

// GetByUsername returns a user by user name.
func (r *Repository) GetByUsername(name string) (*model.User, error) {
	var user model.User
	result := r.db.Where("username = ?", name).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return &user, nil
}

// Create the given user with a registration key.
func (r *Repository) Create(user *model.User, key *model.RegistrationKey) (*model.User, error) {

	if err := r.db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, ErrDuplicate
		}
		return nil, ErrNotCreated
	}

	key.User = *user
	key.UserID = user.ID

	if _, err := r.keyRepo.Create(key); err != nil {
		return nil, ErrNotCreated
	}

	return user, nil
}

// Save a given user to the database.
func (r *Repository) Save(user *model.User) (*model.User, error) {
	result := r.db.Save(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	user = result.Value.(*model.User)
	return user, nil
}

// Activate a user given a key's ID.
func (r *Repository) Activate(id uuid.UUID) error {
	// @FIXME: Activating a user happens by passing in the registration key uuid
	// so the user object has to essentially have the registration object in hand
	// but the user object doesn't have a reference to the registration object, the
	// relationship is that a registration key belongs to a user. It seems
	// it would make more sense to have a registration key activate the user.

	err := r.db.Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		if errors.Is(err, registrationkeys.ErrNotFound) {
			return ErrNotFound
		}
		return ErrStorage
	}
	return nil
}
