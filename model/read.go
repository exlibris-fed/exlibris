package model

import (
	"fmt"
	"net/url"
	"time"

	"github.com/exlibris-fed/gormuuid"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// Read is a many to many model describing a user who read a book. Because GORM does weird things with foreign keys we need to do it manually, unfortunately.
type Read struct {
	gormuuid.UUID
	Created time.Time
	Updated time.Time
	Deleted time.Time
	BookID  string `gorm:"not null"`
	Book    *Book  `gorm:"-"`
	UserID  string `gorm:"not null"`
	User    *User  `gorm:"-"`
}

// ToType returns a representation of a read activity as an ActivityPub object.
func (r *Read) ToType() vocab.Type {
	read := streams.NewActivityStreamsRead()

	if r.User != nil {
		u, err := url.Parse(fmt.Sprintf("https://%s/%s", r.User.ID, r.ID))
		if err == nil {
			id := streams.NewJSONLDIdProperty()
			id.SetIRI(u)
			read.SetJSONLDId(id)
		}

		actor := streams.NewActivityStreamsActorProperty()
		actor.AppendActivityStreamsPerson(r.User.ToType().(vocab.ActivityStreamsPerson))
		read.SetActivityStreamsActor(actor)
	}

	if r.Book != nil {
		document := streams.NewActivityStreamsObjectProperty()
		document.AppendActivityStreamsDocument(r.Book.ToType().(vocab.ActivityStreamsDocument))
		read.SetActivityStreamsObject(document)
	}

	return read
}
