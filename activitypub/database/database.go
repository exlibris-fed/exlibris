// Package database implements the go-fed/activity/Database interface.
package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/exlibris-fed/exlibris/config"
	"github.com/exlibris-fed/exlibris/infrastructure/inbox"
	"github.com/exlibris-fed/exlibris/infrastructure/outbox"
	"github.com/exlibris-fed/exlibris/infrastructure/reads"
	"github.com/exlibris-fed/exlibris/infrastructure/users"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	regexpID        = regexp.MustCompile("/user/([^\\/]+)$")
	regexpOutbox    = regexp.MustCompile("/user/([^\\/]+)/outbox$")
	regexpInbox     = regexp.MustCompile("/user/([^\\/]+)/inbox$")
	regexpRead      = regexp.MustCompile("/user/([^\\/]+)/read/([a-z0-9-]+)$")
	regexpFollowers = regexp.MustCompile("/user/([^\\/]+)/followers$")
)

const (
	// ResultsPerPage is how many results to return in a response
	ResultsPerPage = 10

	// StartPage is the number of pages to skip before returning results
	StartPage = 0
)

// A Database is a connection to a database. It uses the gorm connection, so that we can still use the models.
type Database struct {
	baseURL    string
	cfg        *config.Config
	outboxRepo *outbox.Repository
	inboxRepo  *inbox.Repository
	usersRepo  *users.Repository
	readsRepo  *reads.Repository
	locks      map[*url.URL]*sync.Mutex
}

// New returns a new database object.
func New(db *gorm.DB, cfg *config.Config) *Database {
	uri := url.URL{
		Scheme: cfg.Scheme,
		Host:   cfg.Domain,
	}
	return &Database{
		baseURL:    uri.String(),
		cfg:        cfg,
		outboxRepo: outbox.New(db),
		inboxRepo:  inbox.New(db),
		usersRepo:  users.New(db),
		readsRepo:  reads.New(db),
		locks:      make(map[*url.URL]*sync.Mutex),
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
	log.Println("inboxcontains")
	return d.inboxRepo.Contains(inbox, id)
}

// GetInbox returns the first ordered collection page of the inbox at
// the specified IRI, for prepending new items.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("getting inbox: %s", inboxIRI.String())

	inbox = streams.NewActivityStreamsOrderedCollectionPage()
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(inboxIRI)
	inbox.SetJSONLDId(id)

	// TODO pagination
	entries, err := d.inboxRepo.GetByIRI(inboxIRI)
	if err != nil {
		return nil, err
	}

	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	for _, e := range entries {
		iri, err := url.Parse(e.URI)
		if err != nil {
			log.Printf("error parsing url %s: %s", e.URI, err.Error())
			continue
		}
		orderedItems.AppendIRI(iri)
	}
	inbox.SetActivityStreamsOrderedItems(orderedItems)

	return
}

// SetInbox saves the inbox value given from GetInbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) SetInbox(c context.Context, inbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("set inbox in db")
	items := inbox.GetActivityStreamsOrderedItems()

	inboxIRI := inbox.GetJSONLDId().GetIRI()
	pieces := regexpInbox.FindStringSubmatch(inboxIRI.String())
	if len(pieces) != 2 {
		return fmt.Errorf("invalid url %v", inboxIRI)
	}
	user, err := d.usersRepo.GetByUsername(pieces[1])
	if err != nil {
		return err
	}

	existing := make(map[string]bool)
	existingItems, err := d.inboxRepo.GetByIRI(inboxIRI)
	if err != nil {
		return err
	}
	for _, item := range existingItems {
		existing[item.URI] = true
	}
	for item := items.Begin(); item != nil; item = item.Next() {
		iri := item.GetIRI().String()
		if _, exists := existing[iri]; !exists {
			if err := d.inboxRepo.Create(&model.InboxEntry{
				Base: model.Base{
					ID: uuid.New(),
				},
				User:     *user,
				InboxIRI: inboxIRI.String(),
				URI:      iri,
			}); err != nil {
				return err
			}
		}
	}

	// TODO remove what doesn't exist (requires setting in get)
	return nil
}

// Owns determines if the ActivityPub id is owned by this server.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Owns(c context.Context, id *url.URL) (owns bool, err error) {
	// TODO more robust than this
	owns = id.Scheme == d.cfg.Scheme && id.Host == d.cfg.Domain
	return
}

// ActorForOutbox fetches the actor's IRI for the given outbox IRI.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) ActorForOutbox(c context.Context, outboxIRI *url.URL) (actorIRI *url.URL, err error) {
	log.Println("actorforoutbox")
	actorIRI, err = url.Parse(strings.Replace(outboxIRI.String(), "/outbox", "", 1))
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
	log.Println("in exists, looking at", id.String())
	// TODO things other than reads
	read, getErr := d.readsRepo.Get(id.String())
	if getErr == nil && read.ID == id.String() {
		log.Println("they do exist")
		exists = true
	}
	log.Println("doesnt exist")
	return
}

// Get returns the database entry for the specified id.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) Get(c context.Context, id *url.URL) (value vocab.Type, err error) {
	pieces := regexpRead.FindStringSubmatch(id.String())
	if len(pieces) == 3 {
		return d.getRead(id.String())
	}

	pieces = regexpFollowers.FindStringSubmatch(id.String())
	if len(pieces) == 2 {
		return d.getFollowers(pieces[1])
	}

	pieces = regexpID.FindStringSubmatch(id.String())
	if len(pieces) == 2 {
		return d.getProfile(pieces[1])
	}

	err = fmt.Errorf("don't know how to process uri %v", id)
	return
}

func (d *Database) getRead(strID string) (value vocab.Type, err error) {
	r, err := d.readsRepo.Get(strID)
	if err != nil {
		return
	}
	value = r.ToType()
	return
}

func (d *Database) getFollowers(strID string) (value vocab.Type, err error) {
	u, err := d.usersRepo.GetByUsernameWithFollowers(strID)
	if err != nil {
		return
	}
	value = u.FollowersToType()
	return
}

func (d *Database) getProfile(strID string) (value vocab.Type, err error) {
	u, err := d.usersRepo.GetByUsername(strID)
	if err != nil {
		return
	}
	value = u.ToType()
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
	log.Println("creating!")

	id := asType.GetJSONLDId().GetIRI()

	exists, err := d.Exists(c, id)
	if err != nil {
		return err
	}
	if exists {
		log.Println("exists, not creating")
		return nil
	}

	if asRead, ok := asType.(vocab.ActivityStreamsRead); ok {
		log.Println(asRead)
		//r := new(model.Read)
		// TODO
		// map[@context:https://www.w3.org/ns/activitystreams actor:map[id:localhost:8080/@dconley inbox:http://localhost:8080/user/dconley/inbox name:dconley outbox:http://localhost:8080/user/dconley/outbox preferredUsername:dconley type:Person] id:http://localhost:8080/@dconley/read/75fa9f8e-35a4-4cce-a56f-6b631857a425 object:map[attributedTo:map[id:https://openlibrary.org/authors//authors/OL19430A/ name:Neal Stephenson type:Person] id:https://openlibrary.org/works//works/OL14911626W/ name:Anathem type:Document] to:[http://localhost:8080/user/dconley/followers https://www.w3.org/ns/activitystreams#Public] type:Read]

		return fmt.Errorf("whoops not yet")
		//return d.readsRepo.Create(r)
	}

	// TODO other types

	return fmt.Errorf("dont know how to handle")
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
func (d *Database) GetOutbox(c context.Context, outboxIRI *url.URL) (outbox vocab.ActivityStreamsOrderedCollectionPage, err error) {
	log.Printf("getting outbox: %s", outboxIRI.String())

	outbox = streams.NewActivityStreamsOrderedCollectionPage()
	id := streams.NewJSONLDIdProperty()
	id.SetIRI(outboxIRI)
	outbox.SetJSONLDId(id)

	// TODO pagination
	entries, err := d.outboxRepo.GetByIRI(outboxIRI)
	if err != nil {
		return nil, err
	}

	orderedItems := streams.NewActivityStreamsOrderedItemsProperty()
	for _, e := range entries {
		iri, err := url.Parse(e.URI)
		if err != nil {
			log.Printf("error parsing url %s: %s", e.URI, err.Error())
			continue
		}
		orderedItems.AppendIRI(iri)
	}
	outbox.SetActivityStreamsOrderedItems(orderedItems)

	return
}

// SetOutbox saves the outbox value given from GetOutbox, with new items
// prepended. Note that the new items must not be added as independent
// database entries. Separate calls to Create will do that.
//
// The library makes this call only after acquiring a lock first.
func (d *Database) SetOutbox(c context.Context, outbox vocab.ActivityStreamsOrderedCollectionPage) error {
	log.Println("set outbox in db")
	items := outbox.GetActivityStreamsOrderedItems()

	outboxIRI := outbox.GetJSONLDId().GetIRI()
	pieces := regexpOutbox.FindStringSubmatch(outboxIRI.String())
	if len(pieces) != 2 {
		return fmt.Errorf("invalid url %v", outboxIRI)
	}
	user, err := d.usersRepo.GetByUsername(pieces[1])
	if err != nil {
		return err
	}

	existing := make(map[string]bool)
	existingItems, err := d.outboxRepo.GetByIRI(outboxIRI)
	if err != nil {
		return err
	}
	for _, item := range existingItems {
		existing[item.URI] = true
	}
	for item := items.Begin(); item != nil; item = item.Next() {
		iri := item.GetIRI().String()
		if _, exists := existing[iri]; !exists {
			if err := d.outboxRepo.Create(&model.OutboxEntry{
				Base: model.Base{
					ID: uuid.New(),
				},
				User:      *user,
				OutboxIRI: outboxIRI.String(),
				URI:       iri,
			}); err != nil {
				return err
			}
		}
	}

	// TODO remove what doesn;t exist (requires setting in get)
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
	user := userI.(*model.User)

	id, err = url.Parse(fmt.Sprintf("%s/@%s/read/%v", d.baseURL, strings.ToLower(user.Username), uuid.New().String()))
	log.Printf("*** URL *** %s", id)

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
