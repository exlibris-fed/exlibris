package middleware

import (
	"github.com/exlibris-fed/exlibris/infrastructure/users"
	"github.com/jinzhu/gorm"
)

type Middleware struct {
	userRepo *users.Repository
}

func New(db *gorm.DB) *Middleware {
	return &Middleware{
		userRepo: users.New(db),
	}
}
