package service

import (
	"strings"

	"github.com/exlibris-fed/exlibris/infrastructure/authors"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/exlibris-fed/openlibrary-go"
	"github.com/jinzhu/gorm"
)

// NewAuthor instantiates an author service
func NewAuthor(db *gorm.DB) *Author {
	return &Author{
		db:          db,
		authorsRepo: authors.New(db),
	}
}

// Author is a type for getting authors from a database or API
type Author struct {
	db          *gorm.DB
	authorsRepo *authors.Repository
}

// Get will fetch an author given an OL ID, returning from the database or fetching from the open library api
func (a *Author) Get(id string) *model.Author {
	author, err := a.authorsRepo.GetByID(id)
	if err != nil {
		// Error finding author in DB
		data := a.fetch(id)
		if data == nil {
			return nil
		}
		author = data
	}

	return author
}

func (a *Author) fetch(id string) *model.Author {
	author, err := openlibrary.GetAuthorByID(strings.Replace(id, "/authors/", "", 1))
	if err != nil {
		return nil
	}
	authorModel := model.NewAuthor(author)
	result := a.db.Create(authorModel)
	return result.Value.(*model.Author)
}
