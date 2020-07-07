package model

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
)

// A Book is something that can be read. Currently this only supports things which are in the Library of Congress API, but eventually it'd be great to support fanfiction and other online-only sources.
type Book struct {
	Base
	OpenLibraryID  string   `gorm:"unique;not null" json:"open_library_id"`
	Title          string   `gorm:"not null;index" json:"title"`
	Published      int      `json:"published,omitempty"`
	ISBN           string   `json:"isbn,omitempty"`
	Authors        []Author `gorm:"many2many:book_authors"`
	BookAuthorsID  uuid.UUID
	Subjects       []Subject `gorm:"many2many:book_subjects"`
	BookSubjectsID uuid.UUID
}

// NewBook returns a new instance of a book
func NewBook(id string, title string, published int, isbn string) *Book {
	return &Book{
		OpenLibraryID: id,
		Title:         title,
		Published:     published,
		ISBN:          isbn,
	}
}

// ToType returns a representation of a book as an ActivityPub object.
func (b *Book) ToType() vocab.Type {
	book := streams.NewActivityStreamsDocument()

	u, err := url.Parse(fmt.Sprintf("https://openlibrary.org/works/%s/", b.OpenLibraryID))
	if err == nil {
		id := streams.NewJSONLDIdProperty()
		id.SetIRI(u)
		book.SetJSONLDId(id)
	}

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString(b.Title)
	book.SetActivityStreamsName(name)

	// this isn't ideal since it will default to 1/1 at 12:00:00 am of the year...?
	if b.Published > 0 {
		published := streams.NewActivityStreamsPublishedProperty()
		date := time.Date(b.Published, time.January, 1, 0, 0, 0, 0, time.UTC)
		published.Set(date)
		book.SetActivityStreamsPublished(published)
	}

	authors := streams.NewActivityStreamsAttributedToProperty()
	for _, a := range b.Authors {
		authors.AppendActivityStreamsPerson(a.ToType().(vocab.ActivityStreamsPerson))
	}
	book.SetActivityStreamsAttributedTo(authors)

	return book
}
