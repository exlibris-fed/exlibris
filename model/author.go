package model

import (
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
)

// An Author is someone who has written a Book.
type Author struct {
	Base
	OpenLibraryID string `gorm:"unique;not null" json:"id"`
	Name          string `json:"name"`
	Books         []Book `gorm:"many2many:book_authors"`
	BookAuthorsID uuid.UUID
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
