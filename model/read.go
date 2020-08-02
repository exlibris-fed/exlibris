package model

import (
	"fmt"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
)

const (
	// ContextKeyRead is the context key to use for the read action
	ContextKeyRead ContextKey = "read"
)

// Read is a many to many model describing a user who read a book. Because GORM does weird things with foreign keys we need to do it manually, unfortunately.
type Read struct {
	Base
	Book   Book `gorm:"foreignkey:OpenLibraryID;association_foreignkey:BookID;association_autoupdate:false"`
	BookID string
	User   User `gorm:"association_autoupdate:false"`
	UserID uuid.UUID
}

// ToType returns a representation of a read activity as an ActivityPub object.
func (r *Read) ToType() vocab.Type {
	read := streams.NewActivityStreamsRead()

	u, err := url.Parse(fmt.Sprintf("https://%s/%s", r.User.ID, r.ID))
	if err == nil {
		id := streams.NewJSONLDIdProperty()
		id.SetIRI(u)
		read.SetJSONLDId(id)
	}

	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendActivityStreamsPerson(r.User.ToType().(vocab.ActivityStreamsPerson))
	read.SetActivityStreamsActor(actor)

	document := streams.NewActivityStreamsObjectProperty()
	document.AppendActivityStreamsDocument(r.Book.ToType().(vocab.ActivityStreamsDocument))
	read.SetActivityStreamsObject(document)

	return read
}
