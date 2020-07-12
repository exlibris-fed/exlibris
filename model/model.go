// Package model contains the models used by exlibris. Each model lives in its own file and should be initialized in main.go so that migrations run.
package model

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
)

var (
	// PublicActivityPubIRI is the IRI that indicates an Activity is meant to be visible for general public consumption.
	PublicActivityPubIRI *url.URL

	profileURL   string
	inboxURL     string
	outboxURL    string
	followersURL string
)

func init() {
	if iri, err := url.Parse(pub.PublicActivityPubIRI); err != nil {
		log.Println("can't generate public ActivityPub public IRI", err.Error())
	} else {
		PublicActivityPubIRI = iri
	}

	// TODO do this through viper or our config or something
	scheme := os.Getenv("SCHEME")
	if scheme == "" {
		scheme = "https"
	}
	domain := os.Getenv("DOMAIN")
	baseURL := scheme + "://" + domain
	profileURL = baseURL + "/@%s"
	inboxURL = baseURL + "/user/%s/inbox"
	outboxURL = baseURL + "/user/%s/outbox"
	followersURL = baseURL + "/user/%s/followers"
}

// A ContextKey is a key used to represent a model in a context
type ContextKey string

// A Federater is a representation of a model as an ActivityPub object.
type Federater interface {
	ToType() vocab.Type
}

// Base set of attributes for a model
type Base struct {
	BaseEvents
	ID uuid.UUID `gorm:"primary_key"`
}

// BaseEvents is a set of attributes for events created, updated, deleted
type BaseEvents struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
