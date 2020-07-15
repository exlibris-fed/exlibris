package infrastructure

import (
	"log"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
)

func New(dsn string) *gorm.DB {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	return db
}

func Migrate(db *gorm.DB) {

	// @TODO: Eliminate migrations here and move to github.com/golang-migrate/migrate or liquibase
	db.AutoMigrate(model.Author{})
	db.AutoMigrate(model.Book{})
	db.AutoMigrate(model.OutboxEntry{})
	db.AutoMigrate(model.InboxEntry{})
	db.AutoMigrate(model.Read{})
	db.AutoMigrate(model.Review{})
	db.AutoMigrate(model.Subject{})
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.Follower{})
	db.AutoMigrate(model.RegistrationKey{})
	db.AutoMigrate(model.Cover{})

	db.Table("book_authors").AddForeignKey("author_open_library_id", "authors(open_library_id)", "CASCADE", "CASCADE")
	db.Table("book_authors").AddForeignKey("book_open_library_id", "books(open_library_id)", "CASCADE", "CASCADE")

	db.Table("book_subjects").AddForeignKey("subject_id", "subjects(id)", "CASCADE", "CASCADE")
	db.Table("book_subjects").AddForeignKey("book_open_library_id", "books(open_library_id)", "CASCADE", "CASCADE")

	db.Model(&model.OutboxEntry{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&model.Read{}).AddForeignKey("book_id", "books(open_library_id)", "CASCADE", "CASCADE")
	db.Model(&model.Read{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&model.Review{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.Review{}).AddForeignKey("book_id", "books(open_library_id)", "CASCADE", "CASCADE")

	db.Model(&model.Cover{}).AddForeignKey("book_id", "books(open_library_id)", "CASCADE", "CASCADE")

	db.Model(&model.RegistrationKey{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

}
