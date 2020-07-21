package registrationkeys

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

var usersRows *sqlmock.Rows
var registrationKeysRows *sqlmock.Rows

func setup() {
	ts := time.Date(2020, 7, 11, 12, 0, 0, 0, time.UTC)
	registrationKeysRows = sqlmock.NewRows([]string{"key", "user_id"}).
		AddRow("6f3034c0-0642-4fd8-a040-80e1ee6efaa4", "b3032140-e824-4b39-9be2-47e99f383f2b")
	usersRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "human_id", "username", "email", "display_name"}).
		AddRow(ts, ts, nil, "b3032140-e824-4b39-9be2-47e99f383f2b", "bob@mainframe", "bob", "bob@mainframe", "guardianBob")
}

func teardown() {
	usersRows = nil
	registrationKeysRows = nil
}

func TestGetByUsername(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\" WHERE (\"user_id\" = $1)") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(registrationKeysRows)

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)

	user := &model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
	}

	key, err := repo.GetByUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByUsername_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\" WHERE (\"user_id\" = $1)") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(gorm.ErrRecordNotFound)

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)

	user := &model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
	}

	key, err := repo.GetByUser(user)
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.Equal(t, ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByUsername_ErrStorage(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\" WHERE (\"user_id\" = $1)") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(fmt.Errorf("oops"))

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	user := &model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
	}

	key, err := repo.GetByUser(user)
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.Equal(t, ErrStorage, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\"  WHERE (key = $1) ORDER BY \"registration_keys\".\"key\" ASC LIMIT 1") + "$").
		WithArgs("6f3034c0-0642-4fd8-a040-80e1ee6efaa4").
		WillReturnRows(registrationKeysRows)

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\"  WHERE \"users\".\"deleted_at\" IS NULL AND ((\"id\" IN ($1))) ORDER BY \"users\".\"id\" ASC") + "$").
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(usersRows)

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	key, err := repo.Get(uuid.MustParse("6f3034c0-0642-4fd8-a040-80e1ee6efaa4"))
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\"  WHERE (key = $1) ORDER BY \"registration_keys\".\"key\" ASC LIMIT 1") + "$").
		WithArgs("6f3034c0-0642-4fd8-a040-80e1ee6efaa4").
		WillReturnError(gorm.ErrRecordNotFound)

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	key, err := repo.Get(uuid.MustParse("6f3034c0-0642-4fd8-a040-80e1ee6efaa4"))
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_ErrStorage(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\" WHERE (key = $1) ORDER BY \"registration_keys\".\"key\" ASC LIMIT 1") + "$").
		WithArgs("6f3034c0-0642-4fd8-a040-80e1ee6efaa4").
		WillReturnError(fmt.Errorf("oops"))

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	key, err := repo.Get(uuid.MustParse("6f3034c0-0642-4fd8-a040-80e1ee6efaa4"))
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"registration_keys\" (\"key\",\"user_id\") VALUES ($1,$2) RETURNING \"registration_keys\".\"key\"")+"$").
		WithArgs("33fd50f8-f74e-495f-8c1c-1b791b184c3a", "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(sqlmock.NewRows([]string{"key"}).AddRow("33fd50f8-f74e-495f-8c1c-1b791b184c3a"))
	mock.ExpectCommit()

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)

	key, err := repo.Create(&model.RegistrationKey{
		Key: uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"),
		User: model.User{
			Base: model.Base{
				ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
			},
			HumanID:     "bob@mainframe",
			Username:    "bob",
			Email:       "bob@mainframe",
			DisplayName: "guardianBob",
		},
		UserID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ErrNotCreated(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"registration_keys\" (\"key\",\"user_id\") VALUES ($1,$2) RETURNING \"registration_keys\".\"key\"")+"$").
		WithArgs("33fd50f8-f74e-495f-8c1c-1b791b184c3a", "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)

	key, err := repo.Create(&model.RegistrationKey{
		Key: uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"),
		User: model.User{
			Base: model.Base{
				ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
			},
			HumanID:     "bob@mainframe",
			Username:    "bob",
			Email:       "bob@mainframe",
			DisplayName: "guardianBob",
		},
		UserID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
	})
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.True(t, errors.Is(err, ErrNotCreated))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	conn, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectExec("^" + regexp.QuoteMeta("DELETE FROM \"registration_keys\"  WHERE \"registration_keys\".\"key\" = $1") + "$").
		WithArgs("33fd50f8-f74e-495f-8c1c-1b791b184c3a").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	db, _ := gorm.Open("postgres", conn)
	repo := New(db)

	err := repo.Delete(&model.RegistrationKey{
		Key: uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"),
		User: model.User{
			Base: model.Base{
				ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
			},
			HumanID:     "bob@mainframe",
			Username:    "bob",
			Email:       "bob@mainframe",
			DisplayName: "guardianBob",
		},
		UserID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}
