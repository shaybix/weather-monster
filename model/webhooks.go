package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Webhook describes a webhook for subscribing to a city's temperatures
type Webhook struct {
	ID primitive.ObjectID `bson:"_id"`
	CityID primitive.ObjectID `bson:"city_id"`
	CallbackURL string `bson:"callback_url"`
}

// NewWebhook describes a new webhook to be created
type NewWebhook struct {
	CityID string
	CallbackURL string
}

// WebhookManager describes a webhook model manager
type WebhookManager struct {
	db *mongo.Database
}

func (w *WebhookManager) EnsureIndexes() {
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"city_id": 1, "callback_url": 1}, Options: opt}

	//_, err := w.db.Collection("webhooks").Indexes().CreateOne(
	//	context.TODO(),
	//	index,
	//	options.CreateIndexes().SetMaxTime(10*time.Second),
	//)


	_, err := w.db.Collection("webhooks").Indexes().CreateOne(
		context.TODO(),
		index,
		options.CreateIndexes().SetMaxTime(10*time.Second),
	)
	if err != nil {
		panic(err)
	}
}

// Create creates a new webhook for a given city
func (w *WebhookManager) Create(nw *NewWebhook) (*Webhook, error) {
	id, err := primitive.ObjectIDFromHex(nw.CityID)
	if err != nil {
		return nil, err
	}

	wh := &Webhook{
		ID:          primitive.NewObjectID(),
		CityID:      id,
		CallbackURL: nw.CallbackURL,
	}

	var city City
	if err := w.db.Collection("cities").FindOne(context.Background(), bson.M{"_id": id}).Decode(&city); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}

		return nil, err
	}

	_, err = w.db.Collection("webhooks").InsertOne(context.Background(), wh)
	if err != nil {
		for _, e := range err.(mongo.WriteException).WriteErrors {
			if e.Code == DuplicateErrorCode {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	return wh, nil
}

// Get gets all webhooks associated with a city
func (w *WebhookManager) Get(cityID string) ([]*Webhook, error) {
	oid, err := primitive.ObjectIDFromHex(cityID)
	if err != nil {
		return nil, err
	}

	var webhooks []*Webhook
	cursor, err := w.db.Collection("webhooks").Find(context.Background(), bson.M{"city_id": oid})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}


	for cursor.Next(context.TODO()) {
		var webhook Webhook
		if err := cursor.Decode(&webhook); err != nil {
			log.Println(err)
			continue
		}

		webhooks = append(webhooks, &webhook)
	}

	if len(webhooks) == 0 {
		return nil, ErrNotFound
	}

	return webhooks, nil
}

// Delete deletes a webhook
func (w *WebhookManager) Delete(id string) (*Webhook, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := w.db.Collection("webhooks").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&webhook); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	dr, err := w.db.Collection("webhooks").DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		if dr != nil && dr.DeletedCount == 0 {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &webhook, nil
}

// NewWebhookManager returns a new WebhookManager
func NewWebhookManager(db *mongo.Database) *WebhookManager {
	wm  := &WebhookManager{db}
	wm.EnsureIndexes()
	return  wm
}
