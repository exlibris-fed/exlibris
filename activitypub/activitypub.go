package activitypub

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/exlibris-fed/exlibris/activitypub/clock"
	"github.com/exlibris-fed/exlibris/activitypub/database"
	"github.com/exlibris-fed/exlibris/key"
	"github.com/jinzhu/gorm"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/httpsig"
	"github.com/gorilla/mux"
)

const (
	// UserAgentString is used to identify exlibris in http requests.
	UserAgentString = "exlibris-fed" // TODO version number
)

type contextKey string

const (
	keyUsername contextKey = "username"
)

type ActivityPub struct {
	db    *database.Database
	clock *clock.Clock
}

func New(db *gorm.DB) *ActivityPub {
	return &ActivityPub{
		db:    database.New(db),
		clock: clock.New(),
	}
}

func (ap *ActivityPub) HandleInbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin inbox")
	actor := pub.NewFederatingActor(
		ap,       // common
		ap,       // federating
		ap.db,    // database
		ap.clock, // clock
	)

	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		// how did this happen? I almost want to make it a 500
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := context.WithValue(context.Background(), keyUsername, username)
	if handled, err := actor.PostInbox(c, w, r); err != nil {
		log.Printf("error handling PostInbox for user %s: %s", username, err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		return
	} else if handled {
		log.Printf("handled PostInbox for user %s", username)
		return
	} else if handled, err = actor.GetInbox(c, w, r); err != nil {
		log.Printf("error handling GetInbox for user %s: %s", username, err)
		w.WriteHeader(http.StatusInternalServerError) // TODO
		// Write to w
		return
	} else if handled {
		log.Printf("handled GetInbox for user %s", username)
		return
	}
	log.Println("else...?")
	// else:
	//
	// Handle non-ActivityPub request, such as serving a webpage.
}

func (ap *ActivityPub) HandleOutbox(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlin outbox")
	actor := pub.NewFederatingActor(
		ap,       // common
		ap,       // federating
		ap.db,    // database
		ap.clock, // clock
	)

	// TODO
	c := context.Background()
	// Populate c with request-specific information
	if handled, err := actor.PostOutbox(c, w, r); err != nil {
		// Write to w
		return
	} else if handled {
		return
	} else if handled, err = actor.GetOutbox(c, w, r); err != nil {
		// Write to w
		return
	} else if handled {
		return
	}
	// else:
	//
	// Handle non-ActivityPub request, such as serving a webpage.
}

// AuthenticateGetInbox delegates the authentication of a GET to an
// inbox.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
//
// If an error is returned, it is passed back to the caller of
// GetInbox. In this case, the implementation must not write a
// response to the ResponseWriter as is expected that the client will
// do so when handling the error. The 'authenticated' is ignored.
//
// If no error is returned, but authentication or authorization fails,
// then authenticated must be false and error nil. It is expected that
// the implementation handles writing to the ResponseWriter in this
// case.
//
// Finally, if the authentication and authorization succeeds, then
// authenticated must be true and error nil. The request will continue
// to be processed.
func (ap *ActivityPub) AuthenticateGetInbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO how to determine if logged in?
	log.Println("AuthenticateGetInbox")
	return c, false, nil
}

// AuthenticateGetOutbox delegates the authentication of a GET to an
// outbox.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
//
// If an error is returned, it is passed back to the caller of
// GetOutbox. In this case, the implementation must not write a
// response to the ResponseWriter as is expected that the client will
// do so when handling the error. The 'authenticated' is ignored.
//
// If no error is returned, but authentication or authorization fails,
// then authenticated must be false and error nil. It is expected that
// the implementation handles writing to the ResponseWriter in this
// case.
//
// Finally, if the authentication and authorization succeeds, then
// authenticated must be true and error nil. The request will continue
// to be processed.
func (ap *ActivityPub) AuthenticateGetOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	return
}

// GetOutbox returns the OrderedCollection inbox of the actor for this
// context. It is up to the implementation to provide the correct
// collection for the kind of authorization given in the request.
//
// AuthenticateGetOutbox will be called prior to this.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
func (ap *ActivityPub) GetOutbox(c context.Context, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	return streams.NewActivityStreamsOrderedCollectionPage(), nil
}

// NewTransport returns a new Transport on behalf of a specific actor.
//
// The actorBoxIRI will be either the inbox or outbox of an actor who is
// attempting to do the dereferencing or delivery. Any authentication
// scheme applied on the request must be based on this actor. The
// request must contain some sort of credential of the user, such as a
// HTTP Signature.
//
// The gofedAgent passed in should be used by the Transport
// implementation in the User-Agent, as well as the application-specific
// user agent string. The gofedAgent will indicate this library's use as
// well as the library's version number.
//
// Any server-wide rate-limiting that needs to occur should happen in a
// Transport implementation. This factory function allows this to be
// created, so peer servers are not DOS'd.
//
// Any retry logic should also be handled by the Transport
// implementation.
//
// Note that the library will not maintain a long-lived pointer to the
// returned Transport so that any private credentials are able to be
// garbage collected.
func (ap *ActivityPub) NewTransport(c context.Context, actorBoxIRI *url.URL, gofedAgent string) (t pub.Transport, err error) {
	// TODO don't use the default implementation

	// TODO get user's PK instead of making a new one each time, jfc
	pk, err := key.New()
	if err != nil {
		log.Println("error generating key: " + err.Error())
	}

	t = pub.NewHttpSigTransport(
		&http.Client{},
		gofedAgent+"/"+UserAgentString,
		ap.clock,
		ap.signer([]string{}), // TODO headers
		ap.signer([]string{}), // TODO headers
		"",                    // TODO THIS NEEDS TO BE A PATH TO A PUBLIC KEY (ie /keys/%s)
		pk,
	)
	return
}

func (ap *ActivityPub) signer(headers []string) httpsig.Signer {
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.RSA_SHA256},
		httpsig.DigestSha256,
		headers,
		httpsig.Authorization,
	)
	if err != nil {
		log.Println("error creating signer: " + err.Error())
	}
	return signer
}

// ----- FederatingProtocol ----- //
// Hook callback after parsing the request body for a federated request
// to the Actor's inbox.
//
// Can be used to set contextual information based on the Activity
// received.
//
// Only called if the Federated Protocol is enabled.
//
// Warning: Neither authentication nor authorization has taken place at
// this time. Doing anything beyond setting contextual information is
// strongly discouraged.
//
// If an error is returned, it is passed back to the caller of
// PostInbox. In this case, the DelegateActor implementation must not
// write a response to the ResponseWriter as is expected that the caller
// to PostInbox will do so when handling the error.
func (ap *ActivityPub) PostInboxRequestBodyHook(c context.Context, r *http.Request, activity pub.Activity) (context.Context, error) {
	// TODO
	return c, nil
}

// AuthenticatePostInbox delegates the authentication of a POST to an
// inbox.
//
// If an error is returned, it is passed back to the caller of
// PostInbox. In this case, the implementation must not write a
// response to the ResponseWriter as is expected that the client will
// do so when handling the error. The 'authenticated' is ignored.
//
// If no error is returned, but authentication or authorization fails,
// then authenticated must be false and error nil. It is expected that
// the implementation handles writing to the ResponseWriter in this
// case.
//
// Finally, if the authentication and authorization succeeds, then
// authenticated must be true and error nil. The request will continue
// to be processed.
func (ap *ActivityPub) AuthenticatePostInbox(c context.Context, w http.ResponseWriter, r *http.Request) (out context.Context, authenticated bool, err error) {
	// TODO
	return
}

// Blocked should determine whether to permit a set of actors given by
// their ids are able to interact with this particular end user due to
// being blocked or other application-specific logic.
//
// If an error is returned, it is passed back to the caller of
// PostInbox.
//
// If no error is returned, but authentication or authorization fails,
// then blocked must be true and error nil. An http.StatusForbidden
// will be written in the wresponse.
//
// Finally, if the authentication and authorization succeeds, then
// blocked must be false and error nil. The request will continue
// to be processed.
func (ap *ActivityPub) Blocked(c context.Context, actorIRIs []*url.URL) (blocked bool, err error) {
	// TODO
	return
}

// Callbacks returns the application logic that handles ActivityStreams
// received from federating peers.
//
// Note that certain types of callbacks will be 'wrapped' with default
// behaviors supported natively by the library. Other callbacks
// compatible with streams.TypeResolver can be specified by 'other'.
//
// For example, setting the 'Create' field in the
// FederatingWrappedCallbacks lets an application dependency inject
// additional behaviors they want to take place, including the default
// behavior supplied by this library. This is guaranteed to be compliant
// with the ActivityPub Social protocol.
//
// To override the default behavior, instead supply the function in
// 'other', which does not guarantee the application will be compliant
// with the ActivityPub Social Protocol.
//
// Applications are not expected to handle every single ActivityStreams
// type and extension. The unhandled ones are passed to DefaultCallback.
func (ap *ActivityPub) Callbacks(c context.Context) (wrapped pub.FederatingWrappedCallbacks, other []interface{}, err error) {
	log.Println("callbacks")
	// TODO
	return
}

// DefaultCallback is called for types that go-fed can deserialize but
// are not handled by the application's callbacks returned in the
// Callbacks method.
//
// Applications are not expected to handle every single ActivityStreams
// type and extension, so the unhandled ones are passed to
// DefaultCallback.
func (ap *ActivityPub) DefaultCallback(c context.Context, activity pub.Activity) error {
	// TODO
	log.Println("default callback")
	return nil
}

// MaxInboxForwardingRecursionDepth determines how deep to search within
// an activity to determine if inbox forwarding needs to occur.
//
// Zero or negative numbers indicate infinite recursion.
func (ap *ActivityPub) MaxInboxForwardingRecursionDepth(c context.Context) int {
	// TODO
	return 1
}

// MaxDeliveryRecursionDepth determines how deep to search within
// collections owned by peers when they are targeted to receive a
// delivery.
//
// Zero or negative numbers indicate infinite recursion.
func (ap *ActivityPub) MaxDeliveryRecursionDepth(c context.Context) int {
	// TODO
	return 1
}

// FilterForwarding allows the implementation to apply business logic
// such as blocks, spam filtering, and so on to a list of potential
// Collections and OrderedCollections of recipients when inbox
// forwarding has been triggered.
//
// The activity is provided as a reference for more intelligent
// logic to be used, but the implementation must not modify it.
func (ap *ActivityPub) FilterForwarding(c context.Context, potentialRecipients []*url.URL, a pub.Activity) (filteredRecipients []*url.URL, err error) {
	// TODO
	return
}

// GetInbox returns the OrderedCollection inbox of the actor for this
// context. It is up to the implementation to provide the correct
// collection for the kind of authorization given in the request.
//
// AuthenticateGetInbox will be called prior to this.
//
// Always called, regardless whether the Federated Protocol or Social
// API is enabled.
func (ap *ActivityPub) GetInbox(c context.Context, r *http.Request) (vocab.ActivityStreamsOrderedCollectionPage, error) {
	// TODO
	log.Println("getting inbox and returning vocab")
	return streams.NewActivityStreamsOrderedCollectionPage(), nil
}
