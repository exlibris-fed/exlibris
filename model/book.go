package model

import (
	"fmt"
	"net/url"
	"time"

	"github.com/exlibris-fed/openlibrary-go"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
)

// A Book is something that can be read. Currently this only supports things which are in the Library of Congress API, but eventually it'd be great to support fanfiction and other online-only sources.
type Book struct {
	BaseEvents
	OpenLibraryID string   `gorm:"primary_key" json:"open_library_id"`
	Title         string   `gorm:"not null;index" json:"title"`
	Published     int      `json:"published,omitempty"`
	ISBN          string   `json:"isbn,omitempty"`
	Authors       []Author `gorm:"many2many:book_authors;null"`
	Description   string   `gorm:"null" json:"description"`
	Covers        []Cover  `gorm:"foreignkey:BookID;association_foreignkey:OpenLibraryID;null" json:"covers"`
}

// NewBook returns a new instance of a book
func NewBook(book openlibrary.Work, editions []openlibrary.Edition, authors []Author) *Book {
	result := &Book{
		OpenLibraryID: book.Key,
		Title:         book.Title,
		Authors:       authors,
		Description:   string(book.Description),
	}

	// @TODO: This is just blindly taking the first edition returns in editions, could be smarter?
	if len(editions) > 0 {
		edition := editions[0]
		if len(book.Covers) > 0 {
			result.Covers = append(result.Covers, Cover{Base: Base{ID: uuid.New()}, URL: book.CoverURL(openlibrary.SizeLarge), Type: string(openlibrary.SizeLarge)})
			result.Covers = append(result.Covers, Cover{Base: Base{ID: uuid.New()}, URL: book.CoverURL(openlibrary.SizeMedium), Type: string(openlibrary.SizeMedium)})
			result.Covers = append(result.Covers, Cover{Base: Base{ID: uuid.New()}, URL: book.CoverURL(openlibrary.SizeSmall), Type: string(openlibrary.SizeSmall)})
		}
		if len(edition.Isbn10) > 0 {
			result.ISBN = edition.Isbn10[0]
		}
		if len(edition.Isbn13) > 0 {
			result.ISBN = edition.Isbn13[0]
		}

		if date, err := time.Parse("January 2, 2006", edition.PublishDate); err == nil {
			// @FIXME: we should store int64 instead of int, currently reducing precision
			result.Published = int(date.Unix())
		}
		if date, err := time.Parse("Jan 2, 2006", edition.PublishDate); err == nil {
			// @FIXME: we should store int64 instead of int, currently reducing precision
			result.Published = int(date.Unix())
		}
		if date, err := time.Parse("Jan 2nd, 2006", edition.PublishDate); err == nil {
			// @FIXME: we should store int64 instead of int, currently reducing precision
			result.Published = int(date.Unix())
		}
	}

	return result
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
