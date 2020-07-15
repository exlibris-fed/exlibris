package service

import (
	"fmt"
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

func (b *Book) Get(id string) (*model.Book, error) {
	var book *model.Book
	var err error
	if book, err = b.bookRepository.GetByID("/works/" + id); err != nil {
		// Error finding book in DB
		data, err := b.fetch(id)
		if err != nil {
			return nil, err
		}
		book = data
	}

	return book, nil
}

func (b *Book) fetch(id string) (*model.Book, error) {
	// fetch book from API
	work, err := openlibrary.GetWorkByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch work: %w", err)
	}

	// Fetch editions to get date published
	editions, err := openlibrary.GetEditionsByID(id)
	// @TODO: Persist editions?
	if err != nil {
		return nil, fmt.Errorf("could not fetch editions of work: %w", err)
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
		return nil, fmt.Errorf("could not save work: %w", err)
	}
	return result, nil
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
