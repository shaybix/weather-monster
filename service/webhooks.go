package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
)

// Webhook describes a webhook that is created
type Webhook struct {
	ID          int64  `json:"id"`
	CityID      int64  `json:"city_id"`
	CallbackURL string `json:"callback_url"`
}

// CreateWebhookHandler describes an endpoint that creates a webhook for a specified city
func (m *Manager) CreateWebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(r.FormValue("city_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nw := &model.NewWebhook{
		CityID:      int64(cid),
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
		ID:          wh.ID,
		CityID:      wh.CityID,
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
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	wh, err := m.WM.Delete(int64(id))
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&Webhook{
		ID:          wh.ID,
		CityID:      wh.CityID,
		CallbackURL: wh.CallbackURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
