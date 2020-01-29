package model

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test_CanCreateTemperature(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		tm := NewTemperatureManager(db)

		nt := &NewTemperature{
			CityID: 1,
			Min:    25,
			Max:    29,
		}

		expectedRows := []string{"ID", "min", "max", "timestamp", "city_id"}
		mock.ExpectQuery("INSERT INTO").WillReturnRows(
			sqlmock.NewRows(expectedRows).AddRow(
				1,
				nt.Min,
				nt.Max,
				time.Now().Unix(),
				nt.CityID,
			),
		)

		temp, err := tm.Create(nt)
		r.NoError(err)
		r.NotNil(temp)
		r.Equal(nt.CityID, temp.CityID)
		r.Equal(nt.Min, nt.Min)
		r.Equal(nt.Max, nt.Max)

	}, t)
}

func Test_CannotCreateTemperatureWithNonExistentCityID(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		tm := NewTemperatureManager(db)
		nt := &NewTemperature{
			CityID: 1,
			Min:    25,
			Max:    29,
		}

		mock.ExpectQuery("INSERT INTO").WillReturnError(ErrNotFound)

		temp, err := tm.Create(nt)
		r.Error(err)
		r.Nil(temp)
		r.Equal(err, ErrNotFound)
	}, t)
}
