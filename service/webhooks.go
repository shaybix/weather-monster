package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
	"net/http"
)

// Webhook describes a webhook that is created
type Webhook struct {
	ID string `json:"id"`
	CityID string `json:"city_id"`
	CallbackURL string `json:"callback_url"`
}

// CreateWebhookHandler describes an endpoint that creates a webhook for a specified city
func (m *Manager) CreateWebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	nw := &model.NewWebhook{
		CityID: r.FormValue("city_id"),
		CallbackURL: r.FormValue("callback_url"),

	}

	wh, err := m.WM.Create(nw)
	if err != nil {
		if err == model.ErrAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&Webhook{
		ID: wh.ID.Hex(),
		CityID: wh.CityID.Hex(),
		CallbackURL: wh.CallbackURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)

	return
}

func (m *Manager) DeleteWebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	wh, err := m.WM.Delete(vars["id"])
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&Webhook{
		ID:          wh.ID.Hex(),
		CityID:      wh.CityID.Hex(),
		CallbackURL: wh.CallbackURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
