package model

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test_CanCreateWebhook(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		wm := NewWebhookManager(db)

		nw := &NewWebhook{
			CityID:      1,
			CallbackURL: "http://callback-url.com/callback",
		}

		expectedRows := []string{"ID", "city_id", "callback_url"}
		mock.ExpectQuery("INSERT INTO webhooks").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, nw.CityID, nw.CallbackURL),
		)

		wh, err := wm.Create(nw)
		r.NoError(err)
		r.NotNil(wh)
		r.Equal(nw.CityID, wh.CityID)
		r.Equal(nw.CallbackURL, wh.CallbackURL)
	}, t)
}

func Test_CannotCreateWebhookThatAlreadyExists(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		wm := NewWebhookManager(db)

		nw := &NewWebhook{
			CityID:      1,
			CallbackURL: "http://callback-url.com/callback",
		}

		mock.ExpectQuery("INSERT INTO webhooks").WillReturnError(ErrAlreadyExists)

		wh, err := wm.Create(nw)
		r.Error(err)
		r.Nil(wh)
		r.Equal(ErrAlreadyExists, err)
	}, t)
}

func Test_CanDeleteWebhook(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)
		wm := NewWebhookManager(db)

		expectedRows := []string{"ID", "city_id", "callback_url"}
		mock.ExpectQuery("DELETE webhooks").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, 1, "example.com/callback"),
		)

		wh, err := wm.Delete(1)
		r.NoError(err)
		r.NotNil(wh)
	}, t)
}

func Test_CannotDeleteWebhookThatDoesNotExist(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)
		wm := NewWebhookManager(db)

		mock.ExpectQuery("DELETE webhooks").WillReturnError(ErrNotFound)
		dw, err := wm.Delete(1)
		r.Nil(dw)
		r.Error(err)
		r.Equal(ErrNotFound, err)
	}, t)
}
