// Package database implements the go-fed/activity/Database interface.
package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sync"

	"github.com/exlibris-fed/exlibris/model"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	regexpID     = regexp.MustCompile("^https://(.+)/([^/]+)$")
	regexpOutbox = regexp.MustCompile("^https://(.+)/outbox$")
	regexpInbox  = regexp.MustCompile("^https://(.+)/inbox$")
)

// A Database is a connection to a database. It uses the gorm connection, so that we can still use the models.
type Database struct {
	DB    *gorm.DB
	locks map[*url.URL]*sync.Mutex
}

// New returns a new database object.
func New(db *gorm.DB) *Database {
	return &Database{
		DB:    db,
		locks: make(map[*url.URL]*sync.Mutex),
	}
}

// Lock takes a lock for the object at the specified id. If an error
// is returned, the lock must not have been taken.
//
// The lock must be able to succeed for an id that does not exist in
// the database. This means acquiring the lock does not guarantee the
// entry exists in the database.
//
// Locks are encouraged to be lightweight and in the Go layer, as some
// processes require tight loops acquiring and releasing locks.
//
// Used to ensure race conditions in multiple requests do not occur.
func (d *Database) Lock(c context.Context, id *url.URL) error {
	lock, ok := d.locks[id]
	if !ok {
		lock = new(sync.Mutex)
		d.locks[id] = lock
	}
	lock.Lock()
	return nil
}

// Unlock makes the lock for the object at the specified id available.
// If an error is returned, the lock must have still been freed.
//
// Used to ensure race conditions in multiple requests do not occur.
func (d *Database) Unlock(c context.Context, id *url.URL) error {
	lock, ok := d.locks[id]
	if !ok {
		return fmt.Errorf("lock does not exist for %s", id.String())
	}
	lock.Unlock()
	return nil
}

// InboxContains returns true if the OrderedCollection at 'inbox'
// contains the specified 'id'.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) InboxContains(c context.Context, inbox, id *url.URL) (contains bool, err error) {
	// TODO
	log.Println("inboxcontains")
	return
}

// GetInbox returns the first ordered collection page of the outbox at
// the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	// TODO
	log.Println("getinbox")
	return
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	// TODO
	log.Println("setinbox")
	return nil
}

// Owns returns true if the database has an entry for the IRI and it
// exists in the database.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	// TODO
	log.Println("owns")
	return
}

// ActorForOutbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	// TODO
	pieces := regexpID.FindStringSubmatch(outboxIRI.String())
	actorIRI, err = url.Parse(fmt.Sprintf("https://%s", pieces[1]))
	return
}

// ActorForInbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
	// TODO
	log.Println("actorforinbox")
	return
}

// OutboxForInbox fetches the corresponding actor's outbox IRI for the
// actor's inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
	// TODO
	log.Println("outboxforinbox")
	return
}

// Exists returns true if the database has an entry for the specified
// id. It may not be owned by this application instance.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Exists(c context.Context, id *url.URL) (exists bool, err error) {
	// TODO
	log.Println("exists")
	return
}

// Get returns the database entry for the specified id.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Get(c context.Context, id *url.URL) (value vocab.Type, err error) {
	// TODO
	pieces := regexpID.FindStringSubmatch(id.String())
	var object model.APObject
	// @TODO code wanted read but never used it, this will preload based on relationship
	d.DB.Preload("Read").First(&object, "id = ?", pieces[2])

	// this could be better
	// @TODO: read was never used
	// var read model.Read
	// d.DB.First(&read, "id = ?", object.FKRead)

	book := streams.NewActivityStreamsRead()
	userIRI, err := url.Parse("https://" + pieces[1])
	if err != nil {
		return
	}
	asActor := streams.NewActivityStreamsActorProperty()
	asActor.AppendIRI(userIRI)
	book.SetActivityStreamsActor(asActor)

	value = book
	log.Printf("book = %+v", book)
	log.Printf("value = %+v", value)

	return
}

// Create adds a new entry to the database which must be able to be
// keyed by its id.
//
// Note that Activity values received from federated peers may also be
// created in the database this way if the Federating Protocol is
// enabled. The client may freely decide to store only the id instead of
// the entire value.
//
// The library makes this call only after acquiring a lock first.
//
// Under certain conditions and network activities, Create may be called
// multiple times for the same ActivityStreams object.
func (d *Database) Create(c context.Context, asType vocab.Type) error {
	// TODO what if this isnt a read?
	jid := asType.GetJSONLDId()
	id, err := jid.Serialize()
	if err != nil {
		return err
	}
	pieces := regexpID.FindStringSubmatch(id.(string))

	u, err := uuid.Parse(pieces[2])
	if err != nil {
		return err
	}
	bytes, err := u.MarshalBinary()
	if err != nil {
		return err
	}

	var i int
	d.DB.Model(model.Read{}).Where("id = ?", bytes).Count(&i)

	// we already have this in the database, don't create it again
	if i == 1 {
		return nil
	}

	readI := c.Value(model.ContextKeyRead)
	if readI == nil {
		return fmt.Errorf("no read context")
	}
	//read := readI.(*model.Read)

	log.Println("create is trying to create!!")
	return fmt.Errorf("not implemented")

	//return result.Error
}

// Update sets an existing entry to the database based on the value's
// id.
//
// Note that Activity values received from federated peers may also be
// updated in the database this way if the Federating Protocol is
// enabled. The client may freely decide to store only the id instead of
// the entire value.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Update(c context.Context, asType vocab.Type) error {
	// TODO
	log.Println("update")
	return nil
}

// Delete removes the entry with the given id.
//
// Delete is only called for federated objects. Deletes from the Social
// Protocol instead call Update to create a Tombstone.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Delete(c context.Context, id *url.URL) error {
	// TODO
	log.Println("delete")
	return nil
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) GetOutbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	pieces := regexpOutbox.FindStringSubmatch(inboxIRI.String())

	inbox = streams.NewActivityStreamsOrderedCollectionPage()
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(inboxIRI)
	inbox.SetJSONLDId(id)

	var entries []model.OutboxEntry
	d.DB.Find(&entries, "user_id = ?", pieces[1])
	for _, _ = range entries {
		//log.Printf("found serialized %+v\n", e)
	}

	return
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) SetOutbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	items := inbox.GetActivityStreamsOrderedItems()
	if items == nil {
		return fmt.Errorf("ordered items is nil. is this intended?")
	}
	id, err := inbox.GetJSONLDId().Serialize()
	if err != nil {
		return err
	}
	pieces := regexpOutbox.FindStringSubmatch(id.(string))

	for item := items.Begin(); item != items.End(); item = item.Next() {
		id, err := uuid.Parse(pieces[1])
		if err != nil {
			return err
		}
		// TODO can you try to set things you didn't write? probably!
		resp := d.DB.Create(&model.OutboxEntry{
			Serialized: item.GetIRI().String(),
			User:       model.User{Base: model.Base{ID: id}},
		})
		if resp.Error != nil {
			return resp.Error
		}
	}
	return nil
}

// NewID creates a new IRI id for the provided activity or object. The
// implementation does not need to set the 'id' property and simply
// needs to determine the value.
//
// The go-fed library will handle setting the 'id' property on the
// activity or object provided with the value returned.
func (d *Database) NewID(c context.Context, t vocab.Type) (id *url.URL, err error) {
	userI := c.Value(model.ContextKeyAuthenticatedUser)
	if userI == nil {
		return nil, fmt.Errorf("no authenticated user in context")
	}
	user := userI.(model.User)

	id, err = url.Parse(fmt.Sprintf("https://%s/read/%v", user.ID, uuid.New().String()))

	return
}

// Followers obtains the Followers Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Followers(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	// TODO

	log.Println("followers")
	return
}

// Following obtains the Following Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Following(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	// TODO
	log.Println("following")
	return
}

// Liked obtains the Liked Collection for an actor with the
// given id.
//
// If modified, the library will then call Update.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Liked(c context.Context, actorIRI *url.URL) (followers vocab.ActivityStreamsCollection, err error) {
	// TODO
	log.Println("liked")
	return
}
