package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CannotHandleGetForecastRequestWithNonExistentTemperatures(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		sm := NewServiceManager(db)
		tc := createTestCity(db)

		r := mux.NewRouter()
		r.HandleFunc("/forecasts/{id}", sm.GetForecastHandler).Methods("GET")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/forecasts/%s", ts.URL, tc.ID.Hex())

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status ok got %v", resp.StatusCode)
		}
	}, t)
}

func Test_CannotHandleGetForecastRequestWithNonExistentCity(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/forecasts/{id}", sm.GetForecastHandler).Methods("GET")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/forecasts/%s", ts.URL, primitive.NewObjectID().Hex())

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status not found got %v", resp.StatusCode)
		}
	}, t)
}
