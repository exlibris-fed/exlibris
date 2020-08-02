package outbox

import (
	"net/url"

	"github.com/exlibris-fed/exlibris/model"

	"github.com/jinzhu/gorm"
)

// New created a new outbox repository.
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// A Repository is how you access outboxes persisted in storage.
type Repository struct {
	db *gorm.DB
}

// GetByIRI retrieves the contents of a user's outbox given the IRI to it.
func (r *Repository) GetByIRI(outboxIRI *url.URL) (entries []*model.OutboxEntry, err error) {
	if gormErr := r.db.Where("outbox_iri = ?", outboxIRI.String()).
		Order("created_at desc").
		Find(&entries).
		Error; gormErr != nil {
			err = ErrNotFound
		}
	return
}

// Create persists a new OutboxEntry
func (r *Repository) Create(o *model.OutboxEntry) error {
	if err := r.db.Create(o).Error; err != nil {
		return ErrEntryNotCreated
	}
	return nil
}
