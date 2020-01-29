package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
)

type NewTestCity struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func Test_CanHandleCreateCityRequest(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		nc := &NewTestCity{
			Name:      "Washington",
			Latitude:  62.2324,
			Longitude: 11.3214,
		}

		form := url.Values{}
		form.Add("name", nc.Name)
		form.Add("latitude", fmt.Sprintf("%f", nc.Latitude))
		form.Add("longitude", fmt.Sprintf("%f", nc.Longitude))

		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/cities", sm.CreateCityHandler).Methods("POST")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/cities", ts.URL)

		req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		expectedRows := []string{"ID", "name", "latitude", "longitude", "version"}
		mock.ExpectQuery("INSERT").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, "berlin", 23.232, 34.2323, "random-version-string"),
		)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected Status CREATED, got %v", resp.StatusCode)
		}

	}, t)
}

func Test_CannotHandleCreateCityRequestWithEmptyRequestBody(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/cities", sm.CreateCityHandler).Methods("POST")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/cities", ts.URL)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected Status BAD REQUEST (400), got %v", resp.StatusCode)
		}

	}, t)
}

func Test_CanHandleUpdateCityRequest(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/cities/{id}", sm.UpdateCityHandler).Methods("PATCH")

		ts := httptest.NewServer(r)
		defer ts.Close()

		newName := "Berlin City"

		form := url.Values{}
		form.Add("name", newName)
		form.Add("latitude", fmt.Sprintf("%f", 34.241))
		form.Add("longitude", fmt.Sprintf("%f", 32.3421))
		form.Add("version", "random-version-string")

		url := fmt.Sprintf("%s/cities/%s", ts.URL, "1")

		req, err := http.NewRequest("PATCH", url, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		expectedRows := []string{"ID", "name", "latitude", "longitude", "version"}
		mock.ExpectQuery("UPDATE").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, "berlin", 23.323, 34.1231, "random-versioin-string"),
		)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not dial out: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected http status ok but got %v", resp.StatusCode)
		}

	}, t)
}

func Test_CannotHandleUpdateCityRequestWithEmptyRequestBody(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/cities/{id}", sm.UpdateCityHandler).Methods("PATCH")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/cities/%s", ts.URL, "1")

		req, err := http.NewRequest("PATCH", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not dial out: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected http status ok but got %v", resp.StatusCode)
		}
	}, t)
}

func Test_CanHandleDeleteCityRequest(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/cities/{id}", sm.DeleteCityHandler).Methods("DELETE")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/cities/%s", ts.URL, "1")

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		expectedRows := []string{"ID", "name", "latitude", "longitude", "version"}
		mock.ExpectQuery("DELETE").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, "berlin", 34.131, 31.31312, "random-version-string"),
		)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not dial out: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected http status ok but got %v", resp.StatusCode)
		}

	}, t)

}

func Test_CannotHandleDeleteCityRequestWithNonExistentID(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		sm := NewServiceManager(db)

		r := mux.NewRouter()
		r.HandleFunc("/cities/{id}", sm.DeleteCityHandler).Methods("DELETE")

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/cities/%s", ts.URL, "1")

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Fatalf("could not create delete request: %v", err)
		}

		client := &http.Client{}

		mock.ExpectQuery("DELETE").WillReturnError(model.ErrNotFound)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make delete request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status not found error, got %v", resp.StatusCode)
		}

	}, t)
}
