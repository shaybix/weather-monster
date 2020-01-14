package model

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_CanCreateCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		nc := &NewCity{
			Name:      "NewCity",
			Latitude:  12.4532,
			Longitude: 42.4532,
		}

		city, err := cm.Create(nc)
		r.NoError(err)
		r.NotNil(city)
		r.Equal(nc.Name, city.Name)
		r.Equal(nc.Latitude, city.Latitude)
		r.Equal(nc.Longitude, city.Longitude)
		r.NotNil(city.ID)
		r.NotNil(city.Version)
	}, t)
}

func Test_CannotCreateCityThatExists(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		tc := createTestCity(db)

		cm := NewCityManager(db)

		nc := &NewCity{
			Name:      tc.Name,
			Latitude:  12.4532,
			Longitude: 42.4532,
		}

		city, err := cm.Create(nc)
		r.Nil(city)
		r.Error(err)
		r.Equal(err, ErrAlreadyExists)
	},t )
}

func Test_CanUpdateCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		tc := createTestCity(db)

		updatedName := "changed_city"
		cu := &CityUpdate{
			ID:        tc.ID.Hex(),
			Name:      updatedName,
			Latitude:  tc.Latitude,
			Longitude: tc.Longitude,
			Version:   tc.Version.Hex(),
		}

		cm := NewCityManager(db)

		city, err := cm.Update(cu)
		r.NoError(err)
		r.NotNil(city)
		r.Equal(updatedName, city.Name)
	}, t)
}

func Test_CannotUpdatedNonExistentCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		tc := createTestCity(db)
		cu := &CityUpdate{
			ID:        primitive.NewObjectID().Hex(),
			Name:      tc.Name,
			Latitude:  tc.Latitude,
			Longitude: tc.Longitude,
			Version:   tc.Version.Hex(),
		}

		city, err := cm.Update(cu)
		r.Nil(city)
		r.Error(err)
		r.Equal(err, ErrNotFound)
	}, t)
}

func Test_CanDeleteCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		tc := createTestCity(db)
		city, err := cm.Delete(tc.ID.Hex())
		r.NoError(err)
		r.NotNil(city)
		r.Equal(tc.Name, city.Name)
		r.Equal(tc.Latitude, city.Latitude)
		r.Equal(tc.Longitude, city.Longitude)
		r.Equal(tc.ID, city.ID)
		r.Equal(tc.Version, city.Version)
	}, t)
}

func Test_CannotDeleteNonExistentCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		cm := NewCityManager(db)

		city, err := cm.Delete(primitive.NewObjectID().Hex())
		r.Error(err)
		r.Nil(city)
		r.Equal(err, ErrNotFound)
	}, t)
}
