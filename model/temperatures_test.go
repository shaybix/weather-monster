package model

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_CanCreateTemperature(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		tm := NewTemperatureManager(db)

		city := createTestCity(db)

		nt := &NewTemperature{
			CityID: city.ID.Hex(),
			Min:    25,
			Max:    29,
		}

		temp, err := tm.Create(nt)
		r.NoError(err)
		r.NotNil(temp)
		r.Equal(nt.CityID, temp.CityID.Hex())
		r.Equal(nt.Min, temp.Min)
		r.Equal(nt.Max, temp.Max)

	}, t)
}

func Test_CannotCreateTemperatureWithNonExistentCityID(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)

		tm := NewTemperatureManager(db)
		nt := &NewTemperature{
			CityID: primitive.NewObjectID().Hex(),
			Min:    25,
			Max:    29,
		}

		temp, err := tm.Create(nt)
		r.Error(err)
		r.Nil(temp)
		r.Equal(err, ErrNotFound)
	}, t)
}
