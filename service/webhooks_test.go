package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_CanHandleCreateWebhookRequest(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		m  := NewServiceManager(db)
		r := mux.NewRouter()
		r.HandleFunc("/webhooks", m.CreateWebhookHandler).Methods("POST")
		ts := httptest.NewServer(r)

		tc := createTestCity(db)

		form := url.Values{}
		form.Add("city_id", tc.ID.Hex())
		form.Add("callback_url", "http://example.com/temp")

		urls := fmt.Sprintf("%s/webhooks", ts.URL)
		req, err:= http.NewRequest("POST", urls, strings.NewReader(form.Encode()))
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
			t.Errorf("expected Status created, got %v", resp.StatusCode)
		}
	},t)
}

func Test_CannotHandleCreateWebhookRequestWithNonExistentCity(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		m  := NewServiceManager(db)
		r := mux.NewRouter()
		r.HandleFunc("/webhooks", m.CreateWebhookHandler).Methods("POST")
		ts := httptest.NewServer(r)

		form := url.Values{}
		form.Add("city_id", primitive.NewObjectID().Hex())
		form.Add("callback_url", "http://example.com/temp")

		urls := fmt.Sprintf("%s/webhooks", ts.URL)
		req, err:= http.NewRequest("POST", urls, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

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
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		m  := NewServiceManager(db)
		router := mux.NewRouter()
		router.HandleFunc("/webhooks", m.CreateWebhookHandler).Methods("POST")
		ts := httptest.NewServer(router)

		tc := createTestCity(db)
		whm := model.NewWebhookManager(db)
		nw := &model.NewWebhook{
			CityID:      tc.ID.Hex(),
			CallbackURL: "http://example.com/temps",
		}

		r := require.New(t)
		wh, err := whm.Create(nw)
		r.NotNil(wh)
		r.NoError(err)


		form := url.Values{}
		form.Add("city_id", wh.CityID.Hex())
		form.Add("callback_url", wh.CallbackURL)

		urls := fmt.Sprintf("%s/webhooks", ts.URL)
		req, err:= http.NewRequest("POST", urls, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("expected Status conflict, got %v", resp.StatusCode)
		}
	},t)
}

func Test_CanHandleDeleteWebhookRequest(t *testing.T) {
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		m  := NewServiceManager(db)
		router := mux.NewRouter()
		router.HandleFunc("/webhooks/{id}", m.DeleteWebhookHandler).Methods("DELETE")
		ts := httptest.NewServer(router)

		tc := createTestCity(db)
		whm := model.NewWebhookManager(db)
		nw := &model.NewWebhook{
			CityID:      tc.ID.Hex(),
			CallbackURL: "http://example.com/temps",
		}

		r := require.New(t)
		wh, err := whm.Create(nw)
		r.NotNil(wh)
		r.NoError(err)

		urls := fmt.Sprintf("%s/webhooks/%s", ts.URL, wh.ID.Hex())
		req, err := http.NewRequest("DELETE", urls, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

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
	withServiceTestDB(func(db *mongo.Database, t *testing.T) {
		m  := NewServiceManager(db)
		router := mux.NewRouter()
		router.HandleFunc("/webhooks/{id}", m.DeleteWebhookHandler).Methods("DELETE")
		ts := httptest.NewServer(router)


		urls := fmt.Sprintf("%s/webhooks/%s", ts.URL, primitive.NewObjectID().Hex())
		req, err := http.NewRequest("DELETE", urls, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected Status not found, got %v", resp.StatusCode)
		}
	},t)
}
