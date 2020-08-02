package users

import (
	"crypto"
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
	usersRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "human_id", "username", "email", "display_name", "password", "private_key", "summary", "local", "verified"}).
		AddRow(ts, ts, nil, "b3032140-e824-4b39-9be2-47e99f383f2b", "bob@mainframe", "bob", "bob@mainframe.local", "guardianBob", "password", "", "summary", true, false)
	registrationKeysRows = sqlmock.NewRows([]string{"key", "user_id"}).
		AddRow("6f3034c0-0642-4fd8-a040-80e1ee6efaa4", "b3032140-e824-4b39-9be2-47e99f383f2b")
}

func teardown() {
	usersRows = nil
}

func TestGetByUsername(t *testing.T) {
	setup()
	defer teardown()

	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND ((username = $1)) ORDER BY \"users\".\"id\" ASC LIMIT 1") + "$").
		WithArgs("bob").
		WillReturnRows(usersRows)
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.GetByUsername("bob")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, "bob", user.Username)
	assert.Equal(t, "bob@mainframe.local", user.Email)
	assert.Equal(t, "bob@mainframe", user.HumanID)
	assert.Equal(t, "guardianBob", user.DisplayName)
}

func TestGetByUsername_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND ((username = $1)) ORDER BY \"users\".\"id\" ASC LIMIT 1") + "$").
		WithArgs("bob").
		WillReturnError(gorm.ErrRecordNotFound)
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.GetByUsername("bob")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByUsername_ErrStorage(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND ((username = $1)) ORDER BY \"users\".\"id\" ASC LIMIT 1") + "$").
		WithArgs("bob").
		WillReturnError(fmt.Errorf("oops"))
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.GetByUsername("bob")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	setup()
	defer teardown()

	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"human_id\",\"username\",\"display_name\",\"email\",\"password\",\"private_key\",\"summary\",\"local\",\"verified\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING \"users\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "b3032140-e824-4b39-9be2-47e99f383f2b", "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("b3032140-e824-4b39-9be2-47e99f383f2b"))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("INSERT INTO \"registration_key")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("33fd50f8-f74e-495f-8c1c-1b791b184c3a"))
	mock.ExpectCommit()
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.Create(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	}, &model.RegistrationKey{
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
	assert.NotNil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreate_ErrNotCreated_User(t *testing.T) {
	setup()
	defer teardown()

	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^"+regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"human_id\",\"username\",\"display_name\",\"email\",\"password\",\"private_key\",\"summary\",\"local\",\"verified\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING \"users\".\"id\"")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "b3032140-e824-4b39-9be2-47e99f383f2b", "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true).
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.Create(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	}, &model.RegistrationKey{
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
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrNotCreated))
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreate_ErrNotCreated_Key(t *testing.T) {
	setup()
	defer teardown()

	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"human_id\",\"username\",\"display_name\",\"email\",\"password\",\"private_key\",\"summary\",\"local\",\"verified\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING \"users\".\"id\"") + "$").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("b3032140-e824-4b39-9be2-47e99f383f2b"))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("INSERT INTO \"registration_key")).
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.Create(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	}, &model.RegistrationKey{
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
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrNotCreated))
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestCreate_ErrDuplicate(t *testing.T) {
	setup()
	defer teardown()

	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"id\",\"human_id\",\"username\",\"display_name\",\"email\",\"password\",\"private_key\",\"summary\",\"local\",\"verified\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING \"users\".\"id\"") + "$").
		WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))
	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)

	repo := New(db)
	user, err := repo.Create(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	}, &model.RegistrationKey{
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
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrDuplicate))
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestSave(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	mock.ExpectBegin()
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"users\" SET \"updated_at\" = $1, \"deleted_at\" = $2, \"human_id\" = $3, \"username\" = $4, \"display_name\" = $5, \"email\" = $6, \"password\" = $7, \"private_key\" = $8, \"summary\" = $9, \"local\" = $10, \"verified\" = $11  WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = $12")+"$").
		WithArgs(sqlmock.AnyArg(), nil, "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true, "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	user, err := repo.Save(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_ErrNotFound(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	mock.ExpectBegin()
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"users\" SET")).
		WithArgs(sqlmock.AnyArg(), nil, "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true, "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()
	user, err := repo.Save(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	})

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_ErrStorage(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"users\" SET")).
		WithArgs(sqlmock.AnyArg(), nil, "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true, "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	user, err := repo.Save(&model.User{
		Base: model.Base{
			ID: uuid.MustParse("b3032140-e824-4b39-9be2-47e99f383f2b"),
		},
		HumanID:          "bob@mainframe",
		Username:         "bob",
		Email:            "bob@mainframe.local",
		DisplayName:      "guardianBob",
		CryptoPrivateKey: new(crypto.PrivateKey),
		Local:            true,
		Password:         *new([]byte),
		PrivateKey:       *new([]byte),
		Summary:          "summary",
		Verified:         true,
	})

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivate(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\"")).
		WillReturnRows(registrationKeysRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\"")).
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(usersRows)
	mock.ExpectBegin()
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"users\" SET \"created_at\" = $1, \"updated_at\" = $2, \"deleted_at\" = $3, \"human_id\" = $4, \"username\" = $5, \"display_name\" = $6, \"email\" = $7, \"password\" = $8, \"private_key\" = $9, \"summary\" = $10, \"local\" = $11, \"verified\" = $12  WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = $13")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true, "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("^" + regexp.QuoteMeta("DELETE FROM \"registration_keys\"")).
		WithArgs("6f3034c0-0642-4fd8-a040-80e1ee6efaa4").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectCommit()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	err := repo.Activate(uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"))
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivate_Key_ErrNotFound(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\"")).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	err := repo.Activate(uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivate_User_ErrStorage(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\"")).
		WillReturnRows(registrationKeysRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\"")).
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(usersRows)
	mock.ExpectBegin()
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"users\" SET \"created_at\" = $1, \"updated_at\" = $2, \"deleted_at\" = $3, \"human_id\" = $4, \"username\" = $5, \"display_name\" = $6, \"email\" = $7, \"password\" = $8, \"private_key\" = $9, \"summary\" = $10, \"local\" = $11, \"verified\" = $12  WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = $13")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true, "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	err := repo.Activate(uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivate_RegistrationKeys_ErrStorage(t *testing.T) {
	setup()
	defer teardown()
	conn, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"registration_keys\"")).
		WillReturnRows(registrationKeysRows)
	mock.ExpectQuery("^" + regexp.QuoteMeta("SELECT * FROM \"users\"")).
		WithArgs("b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnRows(usersRows)
	mock.ExpectBegin()
	mock.ExpectExec("^"+regexp.QuoteMeta("UPDATE \"users\" SET \"created_at\" = $1, \"updated_at\" = $2, \"deleted_at\" = $3, \"human_id\" = $4, \"username\" = $5, \"display_name\" = $6, \"email\" = $7, \"password\" = $8, \"private_key\" = $9, \"summary\" = $10, \"local\" = $11, \"verified\" = $12  WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = $13")+"$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "bob@mainframe", "bob", "guardianBob", "bob@mainframe.local", sqlmock.AnyArg(), sqlmock.AnyArg(), "summary", true, true, "b3032140-e824-4b39-9be2-47e99f383f2b").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("^" + regexp.QuoteMeta("DELETE FROM \"registration_keys\"")).
		WithArgs("6f3034c0-0642-4fd8-a040-80e1ee6efaa4").
		WillReturnError(fmt.Errorf("oops"))
	mock.ExpectRollback()
	mock.ExpectRollback()
	db, _ := gorm.Open("postgres", conn)
	repo := New(db)
	err := repo.Activate(uuid.MustParse("33fd50f8-f74e-495f-8c1c-1b791b184c3a"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrStorage))
	assert.NoError(t, mock.ExpectationsWereMet())
}
