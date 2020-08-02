package reads

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
var usersRows *sqlmock.Rows
var readRows *sqlmock.Rows

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
	authorsRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "open_library_id", "name"}).
		AddRow(ts, ts, nil, "/authors/OL7313207A", "Amal El-Mohtar").
		AddRow(ts, ts, nil, "/authors/OL7129451A", "Max Gladstone")
	usersRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "human_id", "username", "email", "display_name"}).
		AddRow(ts, ts, nil, "b3032140-e824-4b39-9be2-47e99f383f2b", "bob@mainframe", "bob", "bob@mainframe", "guardianBob")
	readRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "book_id", "user_id"}).
		AddRow(ts, ts, nil, "/works/OL20473909W", "b3032140-e824-4b39-9be2-47e99f383f2b")
}

func teardown() {
	booksRows = nil
	coversRows = nil
	authorsRows = nil
	bookAuthorsRows = nil
	usersRows = nil
	readRows = nil
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"reads\"  WHERE \"reads\".\"deleted_at\" IS NULL AND ((user_id = $1)) ORDER BY created_at desc") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(readRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\"  WHERE \"books\".\"deleted_at\" IS NULL AND ((\"open_library_id\" IN ($1)))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(booksRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"authors\" INNER JOIN \"book_authors\" ON \"book_authors\".\"author_open_library_id\" = \"authors\".\"open_library_id\" WHERE \"authors\".\"deleted_at\" IS NULL AND ((\"book_authors\".\"book_open_library_id\" IN ($1)))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(bookAuthorsRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"covers\"  WHERE \"covers\".\"deleted_at\" IS NULL AND ((\"book_id\" IN ($1)))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(coversRows)

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	user := &model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:     "bob@mainframe",
		Username:    "bob",
		Email:       "bob@mainframe",
		DisplayName: "guardianBob",
	}
	reads, err := repo.Get(user)
	assert.NoError(t, err)
	assert.NotNil(t, reads)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Greater(t, len(reads[0].Book.Authors), 0, "authors not returned")
	assert.Greater(t, len(reads[0].Book.Covers), 0, "covers not returned")
}

func TestGet_Error(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"reads\"  WHERE \"reads\".\"deleted_at\" IS NULL AND ((user_id = $1)) ORDER BY created_at desc") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(fmt.Errorf("no records found"))

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	user := &model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:     "bob@mainframe",
		Username:    "bob",
		Email:       "bob@mainframe",
		DisplayName: "guardianBob",
	}
	reads, err := repo.Get(user)
	assert.Error(t, err)
	assert.Nil(t, reads)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"reads\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"book_id\",\"user_id\") VALUES ($1,$2,$3,$4,$5,$6) RETURNING \"reads\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "c9f8e9ce-f4bf-4770-88c3-6c5cfe541734", "/works/OL20473909W", "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1f3325e2-ee0d-478f-aecc-122235d7a6ce"))
	mock.ExpectCommit()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	read, err := repo.Create(&model.Read{
		Base: model.Base{
			ID: uuid.MustParse("c9f8e9ce-f4bf-4770-88c3-6c5cfe541734"),
		},
		Book: model.Book{
			OpenLibraryID: "/works/OL20473909W",
		},
		BookID: "/works/OL20473909W",
		User: model.User{
			Base: model.Base{
				ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
			},
		},
		UserID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreate_Error(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"reads\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"book_id\",\"user_id\") VALUES ($1,$2,$3,$4,$5,$6) RETURNING \"reads\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "c9f8e9ce-f4bf-4770-88c3-6c5cfe541734", "/works/OL20473909W", "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(fmt.Errorf("could not update"))
	mock.ExpectRollback()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	read, err := repo.Create(&model.Read{
		Base: model.Base{
			ID: uuid.MustParse("c9f8e9ce-f4bf-4770-88c3-6c5cfe541734"),
		},
		Book: model.Book{
			OpenLibraryID: "/works/OL20473909W",
		},
		BookID: "/works/OL20473909W",
		User: model.User{
			Base: model.Base{
				ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
			},
		},
		UserID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
	})

	assert.Error(t, err)
	assert.Nil(t, read)
	assert.True(t, errors.Is(err, ErrNotCreated))
	assert.NoError(t, mock.ExpectationsWereMet())
}
