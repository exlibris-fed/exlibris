package registrationkeys

import (
	"errors"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound = errors.New("registration key could not be found")
	ErrStorage  = errors.New("error with storage")
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetByUsername(name string) (*model.RegistrationKey, error) {
	var key *model.RegistrationKey
	result := r.db.Table("registration_keys").
		Preload("User").
		Joins("inner join users on registration_keys.user_id = users.id").
		Where("username = ?", name).
		First(&key)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return key, nil
}

func (r *Repository) Get(id uuid.UUID) (*model.RegistrationKey, error) {
	var registrationKey *model.RegistrationKey
	result := r.db.Preload("User").
		Where("key = ?", id).
		First(registrationKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrStorage
	}
	return registrationKey, nil
}

func (r *Repository) Create(key *model.RegistrationKey) (*model.RegistrationKey, error) {
	return nil, nil
}

func (r *Repository) Delete(key *model.RegistrationKey) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(key).Error
	})
}
