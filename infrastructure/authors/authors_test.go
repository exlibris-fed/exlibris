package authors

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	ts := time.Date(2020, 7, 11, 12, 0, 0, 0, time.UTC)
	authorSourceRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "open_library_id", "name"}).
		AddRow(ts, ts, nil, "OL1234567A", "writer mcwriterface")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"authors\" WHERE \"authors\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"authors\".\"open_library_id\" ASC LIMIT 1")).
		WithArgs("OL1234567A").
		WillReturnRows(authorSourceRows)
	bookSourceRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "open_library_id", "title", "published", "isbn", "description"}).
		AddRow(ts, ts, nil, "OL1234567W", "title", 1563408000, 123456789, "description")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"books\" INNER JOIN \"book_authors\" ON \"book_authors\".\"book_open_library_id\" = \"books\".\"open_library_id\" WHERE \"books\".\"deleted_at\" IS NULL AND ((\"book_authors\".\"author_open_library_id\" IN ($1)))")).
		WillReturnRows(bookSourceRows)

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	author, err := repo.GetByID("OL1234567A")
	assert.NoError(t, err)
	assert.NotNil(t, author)
	assert.Equal(t, "writer mcwriterface", author.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"authors\" WHERE \"authors\".\"deleted_at\" IS NULL AND ((open_library_id = $1)) ORDER BY \"authors\".\"open_library_id\" ASC LIMIT 1")).
		WithArgs("OL2345678A").
		WillReturnError(fmt.Errorf("record not found"))

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	author, err := repo.GetByID("OL2345678A")
	assert.Error(t, err)
	assert.Nil(t, author)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	authorSourceRows := sqlmock.NewRows([]string{"open_library_id"}).
		AddRow("OL1234567A")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"authors\" (\"created_at\",\"updated_at\",\"deleted_at\",\"open_library_id\",\"name\") VALUES ($1,$2,$3,$4,$5) RETURNING \"authors\".\"open_library_id\"")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "OL1234567A", "writer mcwriterface").
		WillReturnRows(authorSourceRows)
	mock.ExpectCommit()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	author, err := repo.Create(&model.Author{
		OpenLibraryID: "OL1234567A",
		Name:          "writer mcwriterface",
	})
	assert.NoError(t, err)
	assert.NotNil(t, author)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreate_Error(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"authors\" (\"created_at\",\"updated_at\",\"deleted_at\",\"open_library_id\",\"name\") VALUES ($1,$2,$3,$4,$5) RETURNING \"authors\".\"open_library_id\"")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "OL1234567A", "writer mcwriterface").
		WillReturnError(fmt.Errorf("could not update"))
	mock.ExpectRollback()

	db, _ := gorm.Open("postgres", conn)

	repo := New(db)

	author, err := repo.Create(&model.Author{
		OpenLibraryID: "OL1234567A",
		Name:          "writer mcwriterface",
	})
	assert.Error(t, err)
	assert.Nil(t, author)
	assert.NoError(t, mock.ExpectationsWereMet())

}
