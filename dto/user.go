package dto

const ContextActivityStreams = "https://www.w3.org/ns/activitystreams"
const TypePerson = "Person"

// An ActivityPubUser is a DTO when the request accepts `application/activity+json` (ActivityPub)
type ActivityPubUser struct {
	Context                   []string          `json:"@context"`
	ID                        string            `json:"id"`
	Type                      string            `json:"type"`
	Following                 string            `json:"following"`
	Followers                 string            `json:"followers"`
	Inbox                     string            `json:"inbox"`
	Outbox                    string            `json:"outbox"`
	Username                  string            `json:"preferredUsername"`
	Name                      string            `json:"name"`
	URL                       string            `json:"url"`
	ManuallyApprovesFollowers bool              `json:"manuallyApprovesFollowers"`
	PublicKey                 PublicKey         `json:"publicKey,omitempty"`
	Endpoints                 map[string]string `json:"endpoints"`
	//Icon Object `json:"icon"`
	// featured?
	// summary
	// devices?
	// image (profile header)
}

// A PublicKey is a user's public key
type PublicKey struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
	PEM   string `json:"publicKeyPem"`
}

// An Object is an ActivityPub object.
type Object struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	URL       string `json:"url"`
	// TODO more https://www.w3.org/TR/activitystreams-vocabulary/#dfn-object
}

// NewActivityPubUser returns a struct with default values filled in
func NewActivityPubUser() *ActivityPubUser {
	return &ActivityPubUser{
		Context:                   []string{ContextActivityStreams},
		Type:                      TypePerson,
		ManuallyApprovesFollowers: true, // possible TODO
		Endpoints:                 make(map[string]string),
	}
}
