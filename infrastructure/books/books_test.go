package books

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var booksRows *sqlmock.Rows
var coversRows *sqlmock.Rows
var authorsRows *sqlmock.Rows
var bookAuthorsRows *sqlmock.Rows

func setup() {
	ts := time.Date(2020, 7, 11, 12, 0, 0, 0, time.UTC)
	booksRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "open_library_id", "title", "published", "isbn", "description"}).
		AddRow(ts, ts, nil, "/works/OL20473909W", "This Is How You Lose the Time War", "1563408000", "9781529405231", "description")
	coversRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "type", "url", "book_id"}).
		AddRow(ts, ts, nil, "57fe41db-92ea-40ca-8044-9a698b593d84", "L", "http://covers.openlibrary.org/b/id/9138661-L.jpg", "/works/OL20473909W").
		AddRow(ts, ts, nil, "9f9fc264-bea2-4faf-b3ad-c68dd272bae6", "M", "http://covers.openlibrary.org/b/id/9138661-M.jpg", "/works/OL20473909W").
		AddRow(ts, ts, nil, "848de561-bc1e-4d9c-8e69-84bfb86f1e03", "S", "http://covers.openlibrary.org/b/id/9138661-S.jpg", "/works/OL20473909W")
	bookAuthorsRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "open_library_id", "name", "author_open_library_id", "book_open_library_id"}).
		AddRow(ts, ts, nil, "/authors/OL7313207A", "Amal El-Mohtar", "/authors/OL7313207A", "/works/OL20473909W").
		AddRow(ts, ts, nil, "/authors/OL7129451A", "Max Gladstone", "/authors/OL7129451A", "/works/OL20473909W")
}

func teardown() {
	booksRows = nil
	coversRows = nil
	authorsRows = nil
	bookAuthorsRows = nil
}
func TestGetByID(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\" WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(booksRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"covers\"  WHERE \"covers\".\"deleted_at\" IS NULL AND ((\"book_id\" IN ($1))) ORDER BY \"covers\".\"id\" ASC") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(coversRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"authors\" INNER JOIN \"book_authors\" ON \"book_authors\".\"author_open_library_id\" = \"authors\".\"open_library_id\" WHERE \"authors\".\"deleted_at\" IS NULL AND ((\"book_authors\".\"book_open_library_id\" IN ($1)))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(bookAuthorsRows)

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	book, err := repo.GetByID("/works/OL20473909W")
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, "This Is How You Lose the Time War", book.Title)
	assert.Equal(t, 2, len(book.Authors))
	assert.Equal(t, 3, len(book.Covers))
}

func TestGetByID_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\" WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WillReturnError(gorm.ErrRecordNotFound)

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	book, err := repo.GetByID("/works/OL20473909W")
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_ErrStorage(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\" WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WillReturnError(fmt.Errorf("oops"))

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	book, err := repo.GetByID("/works/OL20473909W")
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	bookSourceRows := sqlmock.NewRows([]string{"open_library_id"}).
		AddRow("/work/OL1234567W")

	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"books\" (\"created_at\",\"updated_at\",\"deleted_at\",\"open_library_id\",\"title\",\"published\",\"isbn\",\"description\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING \"books\".\"open_library_id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "/work/OL1234567W", "title", 123456789, "1234567890", "").
		WillReturnRows(bookSourceRows)
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"authors\" SET \"updated_at\" = $1, \"deleted_at\" = $2, \"name\" = $3  WHERE \"authors\".\"deleted_at\" IS NULL AND \"authors\".\"open_library_id\" = $4")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "writer mcwriterface", "/author/OL1234567A").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^"+regexp.QuoteMeta("INSERT INTO \"book_authors\" (\"book_open_library_id\",\"author_open_library_id\") SELECT $1,$2  WHERE NOT EXISTS (SELECT * FROM \"book_authors\" WHERE \"book_open_library_id\" = $3 AND \"author_open_library_id\" = $4)")+"$").
		WithArgs("/work/OL1234567W", "/author/OL1234567A", "/work/OL1234567W", "/author/OL1234567A").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"covers\" (\"created_at\",\"updated_at\",\"deleted_at\",\"type\",\"url\",\"book_id\") VALUES ($1,$2,$3,$4,$5,$6) RETURNING \"covers\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "L", "https://covers.fake/OL1234567W-L.jpg", "/work/OL1234567W").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(uuid.New()))
	mock.ExpectCommit()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	book, err := repo.Create(&model.Book{
		OpenLibraryID: "/work/OL1234567W",
		Authors: []model.Author{
			{
				Name:          "writer mcwriterface",
				OpenLibraryID: "/author/OL1234567A",
			},
		},
		Covers: []model.Cover{
			{
				Type: "L",
				URL:  "https://covers.fake/OL1234567W-L.jpg",
			},
		},
		Description: "",
		ISBN:        "1234567890",
		Published:   123456789,
		Title:       "title",
	})
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreate_Error(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"books\" (\"created_at\",\"updated_at\",\"deleted_at\",\"open_library_id\",\"title\",\"published\",\"isbn\",\"description\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING \"books\".\"open_library_id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "/work/OL1234567W", "title", 123456789, "1234567890", "").
		WillReturnError(fmt.Errorf("could not update"))
	mock.ExpectRollback()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	book, err := repo.Create(&model.Book{
		OpenLibraryID: "/work/OL1234567W",
		Authors: []model.Author{
			{
				Name:          "writer mcwriterface",
				OpenLibraryID: "/author/OL1234567A",
			},
		},
		Covers: []model.Cover{
			{
				Base: model.Base{
					ID: uuid.New(),
				},
				Type: "L",
				URL:  "https://covers.fake/OL1234567W-L.jpg",
			},
		},
		Description: "",
		ISBN:        "1234567890",
		Published:   123456789,
		Title:       "title",
	})
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.True(t, errors.Is(err, ErrNotCreated))
	assert.NoError(t, mock.ExpectationsWereMet())

}
