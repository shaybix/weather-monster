package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(
		ctx,
		options.Client().
			SetDirect(true).
			ApplyURI("mongodb://db:27017"),
	)
	if err != nil {
		log.Fatalf("error connecting to mongodb: %v", err)
	}
	log.Println("mongodb connection created!")

	db := client.Database("weather")

	mgr := service.NewServiceManager(db)

	r := mux.NewRouter()

	// cities API endpoints
	r.HandleFunc("/cities", mgr.CreateCityHandler).Methods("POST")
	r.HandleFunc("/cities/{id}", mgr.UpdateCityHandler).Methods("PATCH")
	r.HandleFunc("/cities/{id}", mgr.DeleteCityHandler).Methods("DELETE")

	// temperatures API endpoint
	r.HandleFunc("/temperatures", mgr.CreateTemperatureHandler).Methods("POST")

	// forecasts API endpoint
	r.HandleFunc("/forecasts/{id}", mgr.GetForecastHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}

