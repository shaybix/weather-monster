package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Forecast describes a the forecast of a city with the average minimum and maximum temperature in the past 24h
type Forecast struct {
	CityID primitive.ObjectID
	Min    int64
	Max    int64
	Sample int64
}

// ForecastManager describes a forecast model manager
type ForecastManager struct {
	DB *mongo.Database
}

// Get returns the forecast of a city
func (fm *ForecastManager) Get(cid string) (*Forecast, error) {
	cityID, err := primitive.ObjectIDFromHex(cid)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var city City
	if err = fm.DB.Collection("cities").FindOne(ctx, bson.M{"_id": cityID}).Decode(&city); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
	}

	cursor, err := fm.DB.Collection("temperatures").Find(ctx, bson.M{"city_id": cityID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sample int64
	var mins []int64
	var maxs []int64

	forecast := &Forecast{
		CityID: cityID,
	}
	for cursor.Next(context.TODO()) {
		temperature := &Temperature{}
		if err := cursor.Decode(temperature); err != nil {
			log.Println(err)
			continue
		}

		mins = append(mins, temperature.Min)
		maxs = append(maxs, temperature.Max)

		sample++

	}

	if sample == 0 {
		return forecast, nil
	}

	forecast.Sample = sample
	forecast.Min = sum(mins) / int64(len(mins))
	forecast.Min = sum(maxs) / int64(len(maxs))

	return forecast, nil
}

func sum(temps []int64) int64 {
	var total int64
	for _, temp := range temps {
		total = total + temp
	}

	return total
}

// NewForecastManager returns a new ForecastManager
func NewForecastManager(db *mongo.Database) *ForecastManager {
	return &ForecastManager{
		DB: db,
	}
}
