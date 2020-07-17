package model

import (
	"github.com/google/uuid"
)

type RegistrationKey struct {
	Key    uuid.UUID `gorm:"primary_key"`
	User   User      `gorm:"association_autoupdate:false"`
	UserID uuid.UUID
}

func NewRegistrationKey(u User) (*RegistrationKey, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &RegistrationKey{
		Key:    id,
		User:   u,
		UserID: u.ID,
	}, nil
}
