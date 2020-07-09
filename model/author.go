package model

import (
	"fmt"
	"net/url"

	"github.com/exlibris-fed/openlibrary-go"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// An Author is someone who has written a Book.
type Author struct {
	BaseEvents
	OpenLibraryID string `gorm:"primary_key" json:"id"`
	Name          string `json:"name"`
	Books         []Book `gorm:"many2many:book_authors;null"`
}

// ToType returns a representation of an author as an ActivityPub object.
func (a *Author) ToType() vocab.Type {
	author := streams.NewActivityStreamsPerson()

	u, err := url.Parse(fmt.Sprintf("https://openlibrary.org/authors/%s/", a.OpenLibraryID))
	if err == nil {
		id := streams.NewJSONLDIdProperty()
		id.SetIRI(u)
		author.SetJSONLDId(id)
	}

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString(a.Name)
	author.SetActivityStreamsName(name)

	return author
}

// NewAuthor creates an author from an openlibrary Author
func NewAuthor(author openlibrary.Author) *Author {
	return &Author{
		OpenLibraryID: author.Key,
		Name:          author.Name,
	}
}
