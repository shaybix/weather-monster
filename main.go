package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/service"
)

func main() {
	// FIXME: one should be taking the values environment variables
	connStr := "postgresql://postgres:secret@db:5432/weather?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("error closing postgres db: %v", err)
		}
	}()

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

	// webhooks API endpoint
	r.HandleFunc("/webhooks", mgr.CreateWebhookHandler).Methods("POST")
	r.HandleFunc("/webhooks/{id}", mgr.DeleteWebhookHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}
