package model

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func Test_CanCreateWebhook(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)
		tc := createTestCity(db)

		wm := NewWebhookManager(db)

		nw := &NewWebhook{
			CityID:      tc.ID.Hex(),
			CallbackURL: "http://callback-url.com/callback",
		}

		wh, err := wm.Create(nw)
		r.NoError(err)
		r.NotNil(wh)
		r.Equal(nw.CityID, wh.CityID.Hex())
		r.Equal(nw.CallbackURL, wh.CallbackURL)
	}, t)
}

func Test_CannotCreateWebhookForNonExistentCity(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)
		wm := NewWebhookManager(db)

		nw := &NewWebhook{
			CityID:      primitive.NewObjectID().Hex(),
			CallbackURL: "http://callback-url.com/callback",
		}

		wh, err := wm.Create(nw)
		r.Nil(wh)
		r.Error(err)
		r.Equal(ErrNotFound, err)

	}, t)
}

func Test_CannotCreateWebhookThatAlreadyExists(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)
		tc := createTestCity(db)

		wm := NewWebhookManager(db)

		nw := &NewWebhook{
			CityID:      tc.ID.Hex(),
			CallbackURL: "http://callback-url.com/callback",
		}

		wh, err := wm.Create(nw)
		r.NoError(err)
		r.NotNil(wh)
		r.Equal(nw.CityID, wh.CityID.Hex())
		r.Equal(nw.CallbackURL, wh.CallbackURL)

		wh, err = wm.Create(nw)
		r.Nil(wh)
		r.Error(err)
		r.Equal(ErrAlreadyExists, err)
	}, t)
}

func Test_CanDeleteWebhook(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)
		tc := createTestCity(db)
		wm := NewWebhookManager(db)

		nw := &NewWebhook{
			CityID:      tc.ID.Hex(),
			CallbackURL: "http://callback-url.com/callback",
		}

		wh, err := wm.Create(nw)
		r.NoError(err)
		r.NotNil(wh)
		r.Equal(nw.CityID, wh.CityID.Hex())
		r.Equal(nw.CallbackURL, wh.CallbackURL)

		dw, err := wm.Delete(wh.ID.Hex())
		r.NoError(err)
		r.NotNil(dw)
	}, t)
}

func Test_CannotDeleteWebhookThatDoesNotExist(t *testing.T) {
	withTestDB(func(db *mongo.Database, t *testing.T) {
		r := require.New(t)
		wm := NewWebhookManager(db)

		dw, err := wm.Delete(primitive.NewObjectID().Hex())
		r.Nil(dw)
		r.Error(err)
		r.Equal(ErrNotFound, err)
	},t)
}