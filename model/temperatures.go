package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Temperature describes a the temperature of any given day
type Temperature struct {
	ID        primitive.ObjectID `bson:"_id"`
	CityID    primitive.ObjectID `bson:"city_id"`
	Min       int64              `bson:"min"`
	Max       int64              `bson:"max"`
	Timestamp time.Time          `bson:"timestamp"`
}

// NewTemperature describes a new temperature to be added for a city
type NewTemperature struct {
	CityID string
	Min    int64
	Max    int64
}

// TemperatureManager describes a temperature model manager
type TemperatureManager struct {
	DB *mongo.Database
}


// Create creates a temperature entry in the database
func (tm *TemperatureManager) Create(tf *NewTemperature) (*Temperature, error) {
	// TODO: Here the city id is only checked for it being a valid object id,
	// but it should be noted that it should noted, that it is assumed the id given is for a document that exists.
	cid, err := primitive.ObjectIDFromHex(tf.CityID)
	if err != nil {
		return nil, err
	}

	temp := &Temperature{
		ID:        primitive.NewObjectID(),
		CityID:    cid,
		Min:       tf.Min,
		Max:       tf.Max,
		Timestamp: time.Now().UTC(),
	}

	city := &City{}

	if err := tm.DB.Collection("cities").FindOne(context.Background(), bson.M{"_id": cid}).Decode(city); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	_, err = tm.DB.Collection("temperatures").InsertOne(context.Background(), temp)
	if err != nil {
		return nil, err
	}

	return temp, nil
}

// NewTemperatureManager returns a new TemperatureManager
func NewTemperatureManager(db *mongo.Database) *TemperatureManager {
	return  &TemperatureManager{db}
}
