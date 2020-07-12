package activitypub

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/exlibris-fed/exlibris/activitypub/clock"
	"github.com/exlibris-fed/exlibris/activitypub/database"
	"github.com/exlibris-fed/exlibris/config"
	"github.com/exlibris-fed/exlibris/key"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/httpsig"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

const (
	// UserAgentString is used to identify exlibris in http requests.
	UserAgentString = "exlibris-fed" // TODO version number
)

// ActivityPub represents the federating server connection.
type ActivityPub struct {
	db    *database.Database
	clock *clock.Clock
}

// New returns a new ActiityPub object.
func New(db *gorm.DB, cfg *config.Config) *ActivityPub {
	return &ActivityPub{
		db:    database.New(db, cfg),
		clock: clock.New(),
	}
}

func (ap *ActivityPub) NewFederatingActor() pub.FederatingActor {
	return pub.NewFederatingActor(
		ap,       // common
		ap,       // federating
		ap.db,    // database
		ap.clock, // clock
	)
}

func (ap *ActivityPub) NewStreamsHandler() pub.HandlerFunc {
	return pub.NewActivityStreamsHandler(ap.db, ap.clock)
}

// ----- Common ----- //
// AuthenticateGetInbox determines if the request is allowed to access the inbox for a given user.
//
// TODO should it return an error instead of just not authenticated?
func (ap *ActivityPub) AuthenticateGetInbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authenticated bool, err error) {
	log.Println("auth get inbox")
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		return
	}
	user, ok := c.Value(model.ContextKeyAuthenticatedUser).(model.User)
	if !ok {
		// not logged in at all
		return
	}

	// determine if the user is accessing their own
	if strings.ToLower(username) != strings.ToLower(user.Username) {
		return
	}

	// all good, get the inbox
	out = c
	authenticated = true
	return
}

func (ap *ActivityPub) AuthenticateGetOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	log.Println("auth get outbox")
	return
}

func (ap *ActivityPub) GetOutbox(c context.Context, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	log.Println("get outbox")
	return nil, nil
}

func (ap *ActivityPub) NewTransport(c context.Context, actorBoxIRI *url.URL, gofedAgent string) (t pub.Transport, err error) {
	userI := c.Value(model.ContextKeyAuthenticatedUser)
	if userI == nil {
		return nil, fmt.Errorf("no authenticated user in context")
	}
	user, ok := userI.(*model.User)
	if !ok {
		return nil, fmt.Errorf("can't cast to user")
	}

	pk, err := key.DeserializeRSAPrivateKey(user.PrivateKey)
	if err != nil {
		return
	}

	t = pub.NewHttpSigTransport(
		&http.Client{},
		gofedAgent+"/"+UserAgentString,
		ap.clock,
		ap.signer([]string{}), // TODO headers
		ap.signer([]string{}), // TODO headers
		"",                    // TODO THIS NEEDS TO BE A PATH TO A PUBLIC KEY (ie /keys/%s)
		pk.(*rsa.PrivateKey),
	)
	return
}

func (ap *ActivityPub) signer(headers []string) httpsig.Signer {
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.RSA_SHA256},
		httpsig.DigestSha256,
		headers,
		httpsig.Authorization,
		60,
	)
	if err != nil {
		log.Println("error creating signer: " + err.Error())
	}
	return signer
}

// ----- Federating ----- //
func (ap *ActivityPub) PostInboxRequestBodyHook(c context.Context, r *http.Request, activity pub.Activity) (context.Context, error) {
	// TODO
	log.Println("post inbox req body hook")
	return nil, nil
}

func (ap *ActivityPub) AuthenticatePostInbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	log.Println("ath post inbox")
	return
}

func (ap *ActivityPub) Blocked(c context.Context, actorIRIs []*url.URL) (blocked bool, err error) {
	// TODO
	log.Println("blocked")
	return
}

func (ap *ActivityPub) FederatingCallbacks(c context.Context) (wrapped pub.FederatingWrappedCallbacks, other []interface{}, err error) {
	// TODO
	log.Println("fed callbacks")
	return
}

func (ap *ActivityPub) DefaultCallback(c context.Context, activity pub.Activity) error {
	// TODO
	log.Println("default cb")
	return nil
}

func (ap *ActivityPub) MaxInboxForwardingRecursionDepth(c context.Context) int {
	// TODO
	log.Println("max inbox fwd recursion")
	return 25
}

func (ap *ActivityPub) MaxDeliveryRecursionDepth(c context.Context) int {
	// TODO
	depth := 25
	log.Println("using max delivery recursion depth of", depth)
	return depth
}

func (ap *ActivityPub) FilterForwarding(c context.Context, potentialRecipients []*url.URL, a pub.Activity) (filteredRecipients []*url.URL, err error) {
	// TODO
	log.Println("filter forwarding")
	return
}

// GetInbox retrieves a user's inbox. AuthenticateGetInbox already verified that the authenticated user exists and is accessing their own profile.
func (ap *ActivityPub) GetInbox(c context.Context, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	log.Println("get inbox")
	return ap.db.GetInbox(c, r.URL)
}
