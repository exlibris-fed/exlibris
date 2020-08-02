package reviews

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
var reviewsRows *sqlmock.Rows

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
	reviewsRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "book_id", "user_id", "text"}).
		AddRow(ts, ts, nil, "/works/OL20473909W", "b3032140-e824-4b39-9be2-47e99f383f2b", "I had a hard time with this book")
}

func teardown() {
	booksRows = nil
	coversRows = nil
	authorsRows = nil
	bookAuthorsRows = nil
	usersRows = nil
	reviewsRows = nil
}

func TestGetReviews(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"reviews\"  WHERE \"reviews\".\"deleted_at\" IS NULL AND ((book_id = $1))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(reviewsRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\"  WHERE \"users\".\"deleted_at\" IS NULL AND ((\"id\" IN ($1)))") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(usersRows)
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	reviews, err := repo.GetReviews(&model.Book{
		OpenLibraryID: "/works/OL20473909W",
	})

	assert.NoError(t, err)
	assert.NotNil(t, reviews)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetReviews_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"reviews\"  WHERE \"reviews\".\"deleted_at\" IS NULL AND ((book_id = $1))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnError(gorm.ErrRecordNotFound)
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	reviews, err := repo.GetReviews(&model.Book{
		OpenLibraryID: "/works/OL20473909W",
	})

	assert.Error(t, err)
	assert.Nil(t, reviews)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreateReview(t *testing.T) {
	setup()
	defer teardown()

	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\"  WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(booksRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"covers\"  WHERE \"covers\".\"deleted_at\" IS NULL AND ((\"book_id\" IN ($1))) ORDER BY \"covers\".\"id\" ASC") + "$").
		WithArgs(("/works/OL20473909W")).
		WillReturnRows(coversRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"authors\" INNER JOIN \"book_authors\" ON \"book_authors\".\"author_open_library_id\" = \"authors\".\"open_library_id\" WHERE \"authors\".\"deleted_at\" IS NULL AND ((\"book_authors\".\"book_open_library_id\" IN ($1)))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(bookAuthorsRows)
	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"reviews\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"book_id\",\"user_id\",\"text\") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING \"reviews\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "/works/OL20473909W", "b3032140-e824-4b39-9be2-47e99f383f2b", "text").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("10698c21-f094-4a83-8ec7-3221fa9e806e"))
	mock.ExpectCommit()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	review, err := repo.CreateReview(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:     "bob@mainframe",
		Username:    "bob",
		Email:       "bob@mainframe",
		DisplayName: "guardianBob",
	}, &model.Book{
		OpenLibraryID: "/works/OL20473909W",
	}, "text", 5)
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreateReview_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\"  WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnError(gorm.ErrRecordNotFound)

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	review, err := repo.CreateReview(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:     "bob@mainframe",
		Username:    "bob",
		Email:       "bob@mainframe",
		DisplayName: "guardianBob",
	}, &model.Book{
		OpenLibraryID: "/works/OL20473909W",
	}, "text", 5)
	assert.Error(t, err)
	assert.Nil(t, review)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())

}
func TestCreateReview_ErrStorage(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\"  WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnError(fmt.Errorf("oops"))

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	review, err := repo.CreateReview(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:     "bob@mainframe",
		Username:    "bob",
		Email:       "bob@mainframe",
		DisplayName: "guardianBob",
	}, &model.Book{
		OpenLibraryID: "/works/OL20473909W",
	}, "text", 5)
	assert.Error(t, err)
	assert.Nil(t, review)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestCreateReview_ErrNotCreated(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"books\"  WHERE \"books\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"books\".\"open_library_id\" ASC LIMIT 1") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(booksRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"covers\"  WHERE \"covers\".\"deleted_at\" IS NULL AND ((\"book_id\" IN ($1))) ORDER BY \"covers\".\"id\" ASC") + "$").
		WithArgs(("/works/OL20473909W")).
		WillReturnRows(coversRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"authors\" INNER JOIN \"book_authors\" ON \"book_authors\".\"author_open_library_id\" = \"authors\".\"open_library_id\" WHERE \"authors\".\"deleted_at\" IS NULL AND ((\"book_authors\".\"book_open_library_id\" IN ($1)))") + "$").
		WithArgs("/works/OL20473909W").
		WillReturnRows(bookAuthorsRows)
	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"reviews\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"book_id\",\"user_id\",\"text\") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING \"reviews\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "/works/OL20473909W", "b3032140-e824-4b39-9be2-47e99f383f2b", "text").
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	review, err := repo.CreateReview(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:     "bob@mainframe",
		Username:    "bob",
		Email:       "bob@mainframe",
		DisplayName: "guardianBob",
	}, &model.Book{
		OpenLibraryID: "/works/OL20473909W",
	}, "text", 5)
	assert.Error(t, err)
	assert.Nil(t, review)
	assert.True(t, errors.Is(err, ErrNotCreated))
	assert.NoError(t, mock.ExpectationsWereMet())

}
