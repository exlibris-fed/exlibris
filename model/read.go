package model

import (
	"log"
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
	ID string `gorm:"primary_key"`
	BaseEvents
	Book   Book `gorm:"foreignkey:OpenLibraryID;association_foreignkey:BookID"`
	BookID string
	User   User
	UserID uuid.UUID
}

// ToType returns a representation of a read activity as an ActivityPub object.
func (r *Read) ToType() vocab.Type {
	read := streams.NewActivityStreamsRead()

	u, err := url.Parse(r.ID)
	if err != nil {
		log.Printf("error generating url ID for read '%s': %s", r.ID, err.Error())
		return nil
	}
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(u)
	read.SetJSONLDId(id)

	actor := streams.NewActivityStreamsActorProperty()
	actor.AppendActivityStreamsPerson(r.User.ToType().(vocab.ActivityStreamsPerson))
	read.SetActivityStreamsActor(actor)

	document := streams.NewActivityStreamsObjectProperty()
	document.AppendActivityStreamsDocument(r.Book.ToType().(vocab.ActivityStreamsDocument))
	read.SetActivityStreamsObject(document)

	toProperty := streams.NewActivityStreamsToProperty()
	toProperty.AppendIRI(r.User.FollowersIRI())
	if PublicActivityPubIRI != nil {
		toProperty.AppendIRI(PublicActivityPubIRI)
	}
	read.SetActivityStreamsTo(toProperty)

	return read
}
