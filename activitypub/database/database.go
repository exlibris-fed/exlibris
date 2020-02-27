// Package database implements the go-fed/activity/Database interface.
package database

import (
    "context"
    "net/url"

    "github.com/jinzhu/gorm"
    "github.com/go-fed/activity/streams/vocab"
)

// A Database is a connection to a database. It currently uses the gorm connection, which we may want to revisit later?
type Database struct {
    db *gorm.DB
}

// New returns a new database object.
func New(db *gorm.DB) *Database {
    return &Database{
        db: db,
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
    // TODO
    return nil
}

// Unlock makes the lock for the object at the specified id available.
// If an error is returned, the lock must have still been freed.
//
// Used to ensure race conditions in multiple requests do not occur.
func (d *Database) Unlock(c context.Context, id *url.URL) error {
    // TODO
    return nil
}

// InboxContains returns true if the OrderedCollection at 'inbox'
// contains the specified 'id'.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) InboxContains(c context.Context, inbox, id *url.URL) (contains bool, err error) {
    // TODO
    return
}

// GetInbox returns the first ordered collection page of the outbox at
// the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
    // TODO
    return
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
    // TODO
    return nil
}

// Owns returns true if the database has an entry for the IRI and it
// exists in the database.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Owns(c context.Context, id *url.URL) (owns bool, err error) {
    // TODO
    return
}

// ActorForOutbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
    // TODO
    return
}

// ActorForInbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) ActorForInbox(c context.Context, inboxIRI *url.URL) (actorIRI *url.URL, err error) {
    // TODO
    return
}

// OutboxForInbox fetches the corresponding actor's outbox IRI for the
// actor's inbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) OutboxForInbox(c context.Context, inboxIRI *url.URL) (outboxIRI *url.URL, err error) {
    // TODO
    return
}

// Exists returns true if the database has an entry for the specified
// id. It may not be owned by this application instance.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Exists(c context.Context, id *url.URL) (exists bool, err error) {
    // TODO
    return
}

// Get returns the database entry for the specified id.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Get(c context.Context, id *url.URL) (value vocab.Type, err error) {
    // TODO
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
    // TODO
    return nil
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
    return nil
}

// GetOutbox returns the first ordered collection page of the outbox
// at the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) GetOutbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
    // TODO
    return
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) SetOutbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
    // TODO
    return nil
}

// NewId creates a new IRI id for the provided activity or object. The
// implementation does not need to set the 'id' property and simply
// needs to determine the value.
//
// The go-fed library will handle setting the 'id' property on the
// activity or object provided with the value returned.
func (d *Database) NewId(c context.Context, t vocab.Type) (id *url.URL, err error) {
    // TODO
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
    return
}
