package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_CanHandleCreateTemperatureRequest(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		sm := NewServiceManager(db)
		tc := createTestCity(db)

		r := mux.NewRouter()
		r.HandleFunc("/temperatures", sm.CreateTemperatureHandler).Methods("POST")

		ts := httptest.NewServer(r)

		f := url.Values{}
		f.Add("city_id", tc.ID.Hex())
		f.Add("min", "20")
		f.Add("max", "25")

		url := fmt.Sprintf("%s/temperatures", ts.URL)

		req, err := http.NewRequest("POST", url, strings.NewReader(f.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}


		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status created got %v", resp.StatusCode)
		}

	}, t)
}

func Test_CannotHandleCreateTemperatureRequestWithNonExistentCity(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/temperatures", sm.CreateTemperatureHandler).Methods("POST")

		ts := httptest.NewServer(r)

		f := url.Values{}
		f.Add("city_id", primitive.NewObjectID().Hex())
		f.Add("min", "20")
		f.Add("max", "25")

		url := fmt.Sprintf("%s/temperatures", ts.URL)

		req, err := http.NewRequest("POST", url, strings.NewReader(f.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}


		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status internal server error, got %v", resp.StatusCode)
		}

	}, t)
}
