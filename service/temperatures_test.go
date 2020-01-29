package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
)

func Test_CanHandleCreateTemperatureRequest(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/temperatures", sm.CreateTemperatureHandler).Methods("POST")

		ts := httptest.NewServer(r)

		f := url.Values{}
		f.Add("city_id", "1")
		f.Add("min", "20")
		f.Add("max", "25")

		url := fmt.Sprintf("%s/temperatures", ts.URL)

		req, err := http.NewRequest("POST", url, strings.NewReader(f.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		expectedRows := []string{"ID", "min", "max", "city_id", "timestamp"}
		mock.ExpectQuery("INSERT").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, 20, 24, 1, time.Now().Unix()),
		)
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
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/temperatures", sm.CreateTemperatureHandler).Methods("POST")

		ts := httptest.NewServer(r)
		defer ts.Close()

		f := url.Values{}
		f.Add("city_id", "1")
		f.Add("min", "20")
		f.Add("max", "25")

		url := fmt.Sprintf("%s/temperatures", ts.URL)

		req, err := http.NewRequest("POST", url, strings.NewReader(f.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		mock.ExpectQuery("INSERT").WillReturnError(model.ErrNotFound)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status internal server error, got %v", resp.StatusCode)
		}

	}, t)
}
