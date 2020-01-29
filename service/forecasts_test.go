package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
)

func Test_CannotHandleGetForecastRequestWithNonExistentTemperatures(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/forecasts/{id}", sm.GetForecastHandler).Methods("GET")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/forecasts/%s", ts.URL, "1")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		expectedRows := []string{"ID", "min", "max", "timestamp", "city_id"}
		mock.ExpectQuery("SELECT *").WithArgs(1).WillReturnRows(sqlmock.NewRows(expectedRows))

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
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/forecasts/{id}", sm.GetForecastHandler).Methods("GET")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/forecasts/%s", ts.URL, "1")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		mock.ExpectQuery("SELECT").WillReturnError(model.ErrNotFound)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status not found got %v", resp.StatusCode)
		}
	}, t)
}
