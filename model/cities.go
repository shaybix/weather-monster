package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// City describes a city in the world, e.g. Berlin
type City struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Latitude  float64            `bson:"latitude"`
	Longitude float64            `bson:"longitude"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Version   primitive.ObjectID `bson:"version"`
}

// NewCity describes the form values of a new city
type NewCity struct {
	Name      string
	Latitude  float64
	Longitude float64
}

// CityUpdate describes the form values of a city to be updated
type CityUpdate struct {
	ID        string
	Name      string
	Latitude  float64
	Longitude float64
	Version   string
}

// CityManager describes a city model mannager
type CityManager struct {
	db *mongo.Database
}

// EnsureIndexes ensures indexes are created
func (cm *CityManager) EnsureIndexes() {
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"name": 1}, Options: opt}

	_, err := cm.db.Collection("cities").Indexes().CreateOne(
		context.TODO(),
		index,
		options.CreateIndexes().SetMaxTime(10*time.Second),
	)
	if err != nil {
		panic(err)
	}
}

// Create creates a new non-existing entry of a city in the database
func (cm *CityManager) Create(nc *NewCity) (*City, error) {
	city := &City{
		ID:        primitive.NewObjectID(),
		Name:      nc.Name,
		Latitude:  nc.Latitude,
		Longitude: nc.Longitude,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Version:   primitive.NewObjectID(),
	}

	_, err := cm.db.Collection("cities").InsertOne(context.Background(), city)
	if err != nil {
		for _, e := range err.(mongo.WriteException).WriteErrors {
			if e.Code == DuplicateErrorCode {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	return city, nil
}

// Update updates an existing city  in the database
func (cm *CityManager) Update(cu *CityUpdate) (*City, error) {

	id, err := primitive.ObjectIDFromHex(cu.ID)
	if err != nil {
		return nil, err
	}

	ver, err := primitive.ObjectIDFromHex(cu.Version)
	if err != nil {
		return nil, err
	}

	city := &City{}
	if err := cm.db.Collection("cities").FindOne(context.Background(), bson.M{"_id": id, "version": ver}).Decode(city); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	city.Name = cu.Name
	city.Latitude = cu.Latitude
	city.Longitude = cu.Longitude
	city.UpdatedAt = time.Now().UTC()
	city.Version = primitive.NewObjectID()

	_, err = cm.db.Collection("cities").UpdateOne(context.Background(), bson.M{"_id": id, "version": ver}, bson.M{"$set": city})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return city, nil
}

// Delete deletes an existing city in the database
func (cm *CityManager) Delete(id string) (*City, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	city := &City{}
	if err := cm.db.Collection("cities").FindOne(context.Background(), bson.M{"_id": oid}).Decode(city); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	dr, err := cm.db.Collection("cities").DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return nil, err
	} else if dr.DeletedCount == 0 {
		return nil, fmt.Errorf("document with id: %v not found", id)
	}

	return city, nil
}

// NewCityManager returns a new CityManager
func NewCityManager(db *mongo.Database) *CityManager {
	cm := &CityManager{db}
	cm.EnsureIndexes()

	return cm
}
