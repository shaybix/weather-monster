package model

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func withTestDB(f func(db *mongo.Database, t *testing.T), t *testing.T) {
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

func createTestCity(db *mongo.Database) *City {
	city := &City{
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

func createTestTemperatures(db *mongo.Database) (cityID string) {
	city := createTestCity(db)

	temps := []interface{}{
		Temperature{ID: primitive.NewObjectID(), CityID: city.ID, Min: 17, Max: 22, Timestamp: time.Now().Add(time.Minute * 3).UTC()},
		Temperature{ID: primitive.NewObjectID(), CityID: city.ID, Min: 19, Max: 25, Timestamp: time.Now().Add(time.Minute * 5).UTC()},
		Temperature{ID: primitive.NewObjectID(), CityID: city.ID, Min: 22, Max: 28, Timestamp: time.Now().Add(time.Minute * 8).UTC()},
		Temperature{ID: primitive.NewObjectID(), CityID: city.ID, Min: 17, Max: 25, Timestamp: time.Now().Add(time.Minute * 13).UTC()},
		Temperature{ID: primitive.NewObjectID(), CityID: city.ID, Min: 18, Max: 22, Timestamp: time.Now().Add(time.Minute * 17).UTC()},
		Temperature{ID: primitive.NewObjectID(), CityID: city.ID, Min: 19, Max: 24, Timestamp: time.Now().Add(time.Minute * 30).UTC()},
	}

	_, err := db.Collection("temperatures").InsertMany(context.Background(), temps)
	if err != nil {
		panic(err)
	}

	return city.ID.Hex()
}
