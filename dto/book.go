package dto

import (
	"github.com/exlibris-fed/exlibris/model"
)

type Book struct {
	Title     string   `json:"title"`
	Authors   []string `json:"authors"`
	Published int      `json:"published"`
	ISBN      string   `json:"isbn,omitempty"`
	Subjects  []string `json:"subjects"`
	Covers map[string]string `json:"covers"`
}

func newBook(b *model.Book) *Book {
	return &Book{
		Title: b.Title,
	}
}