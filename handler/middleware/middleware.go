package middleware

import (
	"github.com/jinzhu/gorm"
)

type Middleware struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Middleware {
	return &Middleware{
		db: db,
	}
}
