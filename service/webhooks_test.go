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

func Test_CanHandleCreateWebhookRequest(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		m := NewServiceManager(db)
		r := mux.NewRouter()
		r.HandleFunc("/webhooks", m.CreateWebhookHandler).Methods("POST")
		ts := httptest.NewServer(r)

		form := url.Values{}
		form.Add("city_id", "1")
		form.Add("callback_url", "http://example.com/temp")

		urls := fmt.Sprintf("%s/webhooks", ts.URL)
		req, err := http.NewRequest("POST", urls, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		expectedRows := []string{"ID", "city_id", "callback_url"}
		mock.ExpectQuery("INSERT INTO").WillReturnRows(
			sqlmock.NewRows(expectedRows).
				AddRow(1, 1, "example.com/webhook"),
		)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected Status created, got %v", resp.StatusCode)
		}
	}, t)
}

func Test_CannotHandleCreateWebhookRequestWithNonExistentCity(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		m := NewServiceManager(db)
		r := mux.NewRouter()
		r.HandleFunc("/webhooks", m.CreateWebhookHandler).Methods("POST")
		ts := httptest.NewServer(r)

		form := url.Values{}
		form.Add("city_id", "1")
		form.Add("callback_url", "http://example.com/temp")

		urls := fmt.Sprintf("%s/webhooks", ts.URL)
		req, err := http.NewRequest("POST", urls, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		mock.ExpectQuery("INSERT INTO").WillReturnError(model.ErrNotFound)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected Status not found, got %v", resp.StatusCode)
		}
	}, t)
}

func Test_CannotHandleCreateWebhookRequestWithAlreadyExistingWebhook(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		m := NewServiceManager(db)
		router := mux.NewRouter()
		router.HandleFunc("/webhooks", m.CreateWebhookHandler).Methods("POST")
		ts := httptest.NewServer(router)

		form := url.Values{}
		form.Add("city_id", "1")
		form.Add("callback_url", "http://example.com/temps")

		urls := fmt.Sprintf("%s/webhooks", ts.URL)
		req, err := http.NewRequest("POST", urls, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		mock.ExpectQuery("INSERT INTO").WillReturnError(model.ErrAlreadyExists)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected Status conflict, got %v", resp.StatusCode)
		}
	}, t)
}

func Test_CanHandleDeleteWebhookRequest(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		m := NewServiceManager(db)
		router := mux.NewRouter()
		router.HandleFunc("/webhooks/{id}", m.DeleteWebhookHandler).Methods("DELETE")
		ts := httptest.NewServer(router)

		urls := fmt.Sprintf("%s/webhooks/%s", ts.URL, "1")
		req, err := http.NewRequest("DELETE", urls, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		expectedRows := []string{"ID", "city_id", "callback_urlr"}
		mock.ExpectQuery("DELETE").WillReturnRows(
			sqlmock.NewRows(expectedRows).AddRow(
				1,
				1,
				"http://example.com/callback",
			),
		)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected Status ok, got %v", resp.StatusCode)
		}
	}, t)
}

func Test_CannotHandleDeleteWebhookRequestWithNonExistentWebhook(t *testing.T) {
	withServiceTestDB(func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T) {
		m := NewServiceManager(db)
		router := mux.NewRouter()
		router.HandleFunc("/webhooks/{id}", m.DeleteWebhookHandler).Methods("DELETE")
		ts := httptest.NewServer(router)

		urls := fmt.Sprintf("%s/webhooks/%s", ts.URL, "1")
		req, err := http.NewRequest("DELETE", urls, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		mock.ExpectQuery("DELETE").WillReturnError(model.ErrNotFound)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected Status not found, got %v", resp.StatusCode)
		}
	}, t)
}
