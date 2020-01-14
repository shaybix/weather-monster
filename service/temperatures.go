package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shaybix/weather-monster/model"
)

// Temperature describes a temperature of a given city at a specific point in time
type Temperature struct {
	ID     string `json:"id"`
	CityID string `json:"city_id"`
	Min    int64  `json:"min"`
	Max    int64  `json:"max"`
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

	nt := &model.NewTemperature{
		CityID: r.FormValue("city_id"),
		Min:    int64(min),
		Max:    int64(max),
	}

	temp, err := m.TM.Create(nt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&Temperature{
		ID:     temp.ID.Hex(),
		CityID: temp.CityID.Hex(),
		Min:    temp.Min,
		Max:    temp.Max,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
