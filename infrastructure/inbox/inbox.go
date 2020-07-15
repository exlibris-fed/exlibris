package inbox

import (
	"net/url"

	"github.com/exlibris-fed/exlibris/model"

	"github.com/jinzhu/gorm"
)

// New created a new inbox repository.
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// A Repository is how you access inboxes persisted in storage.
type Repository struct {
	db *gorm.DB
}

// GetByIRI retrieves the contents of a user's inbox given the IRI to it.
//
// TODO pagination
func (r *Repository) GetByIRI(inboxIRI *url.URL) (entries []*model.InboxEntry, err error) {
	err = r.db.Where("inbox_iri = ?", inboxIRI.String()).
		Order("created_at desc").
		Find(&entries).
		Error
	return
}

// Contains returns whether an inbox contains a specific IRI.
func (r *Repository) Contains(inboxIRI, iri *url.URL) (contains bool, err error) {
	var count int
	err = r.db.Model(&model.InboxEntry{}).
		Where("inbox_iri = ? AND uri = ?", inboxIRI.String(), iri.String()).
		Count(&count).
		Error

	if err != nil {
		return
	}
	contains = count > 0
	return
}

// Create persists a new InboxEntry
func (r *Repository) Create(i *model.InboxEntry) error {
	return r.db.Create(i).Error
}
