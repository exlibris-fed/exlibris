package model

import (
	"github.com/google/uuid"
)

type RegistrationKey struct {
	Key    uuid.UUID
	User   User
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
