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
	ErrNotFound  = errors.New("user could not be found")
	ErrStorage   = errors.New("error with storage")
	ErrDuplicate = errors.New("user already exists")
)

func New(db *gorm.DB) *Repository {
	return &Repository{
		db:      db,
		keyRepo: registrationkeys.New(db),
	}
}

type Repository struct {
	db      *gorm.DB
	keyRepo *registrationkeys.Repository
}

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
