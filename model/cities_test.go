package model

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test_CanCreateCity(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)
		nc := &NewCity{
			Name:      "NewCity",
			Latitude:  12.4532,
			Longitude: 42.4532,
		}

		expectedRows := []string{"ID", "name", "latitude", "longitude", "version"}
		mock.ExpectQuery("INSERT INTO").WillReturnRows(
			sqlmock.NewRows(expectedRows).AddRow(
				1,
				nc.Name,
				nc.Latitude,
				nc.Longitude,
				"random-version-string",
			),
		)

		city, err := cm.Create(nc)
		r.NoError(err)
		r.NotNil(city)
		r.Equal(nc.Name, city.Name)
		r.Equal(nc.Latitude, city.Latitude)
		r.Equal(nc.Longitude, city.Longitude)
	}, t)
}

func Test_CannotCreateCityThatExists(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		mock.ExpectQuery("INSERT INTO").WillReturnError(ErrAlreadyExists)

		nc := &NewCity{
			Name:      "Berlin",
			Latitude:  12.4532,
			Longitude: 42.4532,
		}

		city, err := cm.Create(nc)
		r.Nil(city)
		r.Error(err)
		r.Equal(err, ErrAlreadyExists)
	}, t)
}

func Test_CanUpdateCity(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		cu := &CityUpdate{
			ID:        1,
			Name:      "NewCity",
			Latitude:  23.131,
			Longitude: 42.12131,
			Version:   "randomstring",
		}

		cm := NewCityManager(db)

		changedName := "updated-city"
		expectedRows := []string{"ID", "name", "latitude", "longitude", "version"}
		mock.ExpectQuery("UPDATE FROM").WillReturnRows(
			sqlmock.NewRows(expectedRows).AddRow(
				cu.ID,
				changedName,
				cu.Latitude,
				cu.Longitude,
				cu.Version,
			),
		)

		city, err := cm.Update(cu)
		r.NoError(err)
		r.NotNil(city)
		r.Equal(cu.ID, city.ID)
		r.Equal(changedName, city.Name)
		r.Equal(cu.Latitude, city.Latitude)
		r.Equal(cu.Version, city.Version)
	}, t)
}

func Test_CannotUpdatedNonExistentCity(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		cu := &CityUpdate{
			ID:        1,
			Name:      "NewCity",
			Latitude:  23.131,
			Longitude: 42.12131,
			Version:   "randomstring",
		}

		mock.ExpectQuery("UPDATE FROM").WillReturnError(ErrNotFound)

		city, err := cm.Update(cu)
		r.Nil(city)
		r.Error(err)
		r.Equal(err, ErrNotFound)
	}, t)
}

func Test_CanDeleteCity(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		expectedRows := []string{"ID", "name", "latitude", "longitude", "version"}
		mock.ExpectQuery("DELETE FROM cities").WillReturnRows(
			sqlmock.NewRows(expectedRows).AddRow(
				1,
				"Berlin",
				13.121431,
				44.3421,
				"random-version-string",
			),
		)

		city, err := cm.Delete(1)
		r.NoError(err)
		r.NotNil(city)
	}, t)
}

func Test_CannotDeleteNonExistentCity(t *testing.T) {
	withTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		mock.ExpectQuery("DELETE FROM cities").WillReturnError(ErrNotFound)

		city, err := cm.Delete(1)
		r.Error(err)
		r.Nil(city)
		r.Equal(err, ErrNotFound)
	}, t)
}
