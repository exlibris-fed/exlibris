package service

import (
	"log"

	"github.com/exlibris-fed/openlibrary-go"
	"github.com/jinzhu/gorm"
)

// NewEditions creates a new instance of edititons
func NewEditions(db *gorm.DB) *Editions {
	return &Editions{
		db: db,
	}
}

// Editions acts as a way to fetch editions of OL works
type Editions struct {
	db *gorm.DB
}

// Get returns a list of editions for a given OLID from OL API
func (e *Editions) Get(id string) []openlibrary.Edition {
	// @TODO: Editions are not stored in db, fetch. Maybe we store these?
	editions, err := openlibrary.GetEditionsByID(id)
	if err != nil {
		log.Println("Could not fetch work editions", id, "got error", err)
		return nil
	}
	return editions
}
