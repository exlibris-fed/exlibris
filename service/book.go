package service

import (
	"log"
	"strings"

	"github.com/exlibris-fed/exlibris/infrastructure/authors"
	"github.com/exlibris-fed/exlibris/infrastructure/books"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/exlibris-fed/openlibrary-go"
	"github.com/jinzhu/gorm"
)

func NewBook(db *gorm.DB) *Book {
	return &Book{
		db:             db,
		bookRepository: books.New(db),
	}
}

type Book struct {
	db             *gorm.DB
	bookRepository *books.Repository
}

func (b *Book) Get(id string) *model.Book {
	var book *model.Book
	var err error
	if book, err = b.bookRepository.GetByID(id); err != nil {
		// Error finding book in DB
		data := b.fetch(id)
		if data == nil {
			return nil
		}
		book = data
	}

	return book
}

func (b *Book) fetch(id string) *model.Book {
	// fetch book from API
	work, err := openlibrary.GetWorkByID(id)
	if err != nil {
		log.Println("Could not fetch work", id, "got error", err)

		return nil
	}

	// Fetch editions to get date published
	editions, err := openlibrary.GetEditionsByID(id)
	// @TODO: Persist editions?
	if err != nil {
		log.Println("Could not fetch work editions", id, "got error", err)
		return nil
	}

	// Gather up the authors
	var authors []model.Author
	for _, author := range work.Authors {
		author := NewAuthor(b.db).Get(author.Author.Key)

		if author == nil {
			continue
		}
		authors = append(authors, *author)
	}

	// Assemble all the data into a book
	book := model.NewBook(work, editions, authors)

	result, err := b.bookRepository.Create(book)
	if err != nil {
		log.Println("Could not insert book into DB:", err)
		return nil
	}
	return result
}

func NewAuthor(db *gorm.DB) *Author {
	return &Author{
		db:          db,
		authorsRepo: authors.New(db),
	}
}

type Author struct {
	db          *gorm.DB
	authorsRepo *authors.Repository
}

func (a *Author) Get(id string) *model.Author {
	author, err := a.authorsRepo.GetByID(id)
	if err != nil {
		// Error finding author in DB
		data := a.fetch(id)
		if data == nil {
			return nil
		}
		author = data
	}

	return author
}

func (a *Author) fetch(id string) *model.Author {
	author, err := openlibrary.GetAuthorByID(strings.Replace(id, "/authors/", "", 1))
	if err != nil {
		log.Println("could not fetch author", err)
		return nil
	}
	authorModel := model.NewAuthor(author)
	result := a.db.Create(authorModel)
	return result.Value.(*model.Author)
}

func NewEditions(db *gorm.DB) *Editions {
	return &Editions{
		db: db,
	}
}

type Editions struct {
	db *gorm.DB
}

func (e *Editions) Get(id string) []openlibrary.Edition {
	// @TODO: Editions are not stored in db, fetch. Maybe we store these?
	editions, err := openlibrary.GetEditionsByID(id)
	if err != nil {
		log.Println("Could not fetch work editions", id, "got error", err)
		return nil
	}
	return editions
}
