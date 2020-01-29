package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/shaybix/weather-monster/model"
)

// Temperature describes a temperature of a given city at a specific point in time
type Temperature struct {
	ID     int64 `json:"id"`
	CityID int64 `json:"city_id"`
	Min    int64 `json:"min"`
	Max    int64 `json:"max"`
}

// CreateTemperatureHandler creates temperature for a specific city
func (m *Manager) CreateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	min, err := strconv.Atoi(r.FormValue("min"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	max, err := strconv.Atoi(r.FormValue("max"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(r.FormValue("city_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nt := &model.NewTemperature{
		CityID: int64(cid),
		Min:    int64(min),
		Max:    int64(max),
	}

	temp, err := m.TM.Create(nt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	whs, err := m.WM.Get(temp.CityID)
	if err != nil {
		log.Println(err)
	}

	t := &Temperature{
		ID:     temp.ID,
		CityID: temp.CityID,
		Min:    temp.Min,
		Max:    temp.Max,
	}

	go m.NotifyWebhooks(whs, t)

	resp, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// NotifyWebhooks notifies all
func (m *Manager) NotifyWebhooks(whs []*model.Webhook, temp *Temperature) {

	client := &http.Client{}

	for _, wh := range whs {
		b, err := json.Marshal(temp)
		if err != nil {
			log.Println(err)
			continue
		}

		req, err := http.NewRequest("POST", wh.CallbackURL, bytes.NewReader(b))
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}

	}
}
