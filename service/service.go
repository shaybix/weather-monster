package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shaybix/weather-monster/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func withServiceTestDB(f func(db *mongo.Database, t *testing.T), t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(
		ctx,
		options.Client().
			SetDirect(true).
			ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		log.Fatalf("error connecting to mongodb: %v", err)
	}
	log.Println("mongodb connection created!")

	uid, err := uuid.NewV4()
	if err != nil {
		log.Panic(err)
		return
	}

	db := client.Database(uid.String())
	defer db.Drop(ctx)
	f(db, t)
}

// Manager ...
type Manager struct {
	CM *model.CityManager
	FM *model.ForecastManager
	TM *model.TemperatureManager
}
func createTestCity(db *mongo.Database) *model.City {


	city := &model.City{
		ID:        primitive.NewObjectID(),
		Name:      "new-city",
		Latitude:  23.4342,
		Longitude: 24.3424,
		Version:   primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
	}

	_, err := db.Collection("cities").InsertOne(context.TODO(), city)
	if err != nil {
		log.Panic(err)
	}
	return city
}

// NewServiceManager ...
func NewServiceManager(db *mongo.Database) *Manager {
	return &Manager{
		CM: model.NewCityManager(db),
		FM: model.NewForecastManager(db),
		TM: model.NewTemperatureManager(db),
	}
}
