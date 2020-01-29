package model

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test_CanGetForecast(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		fm := NewForecastManager(db)
		expectedRows := []string{"ID", "min", "max", "timestamp", "city_id"}
		mock.ExpectQuery("SELECT *").WithArgs(1).WillReturnRows(sqlmock.NewRows(expectedRows))

		fc, err := fm.Get(1)
		r.NoError(err)
		r.NotNil(fc)
	}, t)
}

func Test_CannotGetForecastWithNonExistentCity(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		fm := NewForecastManager(db)
		mock.ExpectQuery("SELECT *").WillReturnError(ErrNotFound)

		fc, err := fm.Get(1)
		r.Nil(fc)
		r.Error(err)
		r.Equal(err, ErrNotFound)
	}, t)
}

func Test_GetForecastWithNoTemperaturesForCityReturnsZeroedValues(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		expectedRows := []string{"ID", "min", "max", "timestamp", "city_id"}
		mock.ExpectQuery("SELECT *").WithArgs(1).WillReturnRows(sqlmock.NewRows(expectedRows))

		fm := NewForecastManager(db)
		fc, err := fm.Get(1)
		r.NoError(err)
		r.NotNil(fc)
	}, t)
}
