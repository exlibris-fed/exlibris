package model

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// An Author is someone who has written a Book.
type Author struct {
	Created time.Time
	Updated time.Time
	Deleted time.Time
	ID      string `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
}

// ToType returns a representation of an author as an ActivityPub object.
func (a *Author) ToType() vocab.Type {
	author := streams.NewActivityStreamsPerson()

	u, err := url.Parse(fmt.Sprintf("https://openlibrary.org/authors/%s/", a.ID))
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
