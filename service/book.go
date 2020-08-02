package service

import (
	"fmt"

	"github.com/exlibris-fed/exlibris/infrastructure/books"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/exlibris-fed/openlibrary-go"
	"github.com/jinzhu/gorm"
)

// NewBook instantiates a book service
func NewBook(db *gorm.DB) *Book {
	return &Book{
		db:             db,
		bookRepository: books.New(db),
	}
}

// Book is a type for getting books from a database or API
type Book struct {
	db             *gorm.DB
	bookRepository *books.Repository
}

// Get will fetch a  work given an OL ID, returning from the database or fetching from the open library api
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
