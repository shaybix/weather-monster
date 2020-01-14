package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_CanGetForecast(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		fm := NewForecastManager(db)

		cityID := createTestTemperatures(db)

		fc, err := fm.Get(cityID)
		r.NoError(err)
		r.NotNil(fc)
	}, t)
}

func Test_CannotGetForecastWithNonExistentCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		fm := NewForecastManager(db)

		fc, err := fm.Get(primitive.NewObjectID().Hex())
		r.Nil(fc)
		r.Error(err)
		r.Equal(err, ErrNotFound)
	}, t)
}

func Test_GetForecastWithNoTemperaturesForCityReturnsZeroedValues(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)
		tc := createTestCity(db)

		fm := NewForecastManager(db)

		fc, err := fm.Get(tc.ID.Hex())
		r.NotNil(fc)
		r.NoError(err)
		r.Equal(fc.Min, int64(0))
		r.Equal(fc.Max, int64(0))
		r.Equal(fc.Sample, int64(0))

	}, t)
}
