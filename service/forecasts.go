package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
	"net/http"
)

// Forecast describes the forecast of a given city in the last 24 hours
type Forecast struct {
	CityID string `json:"city_id"`
	Max int64 `json:"max"`
	Min int64 `json:"min"`
	Sample int64 `json:"sample"`
}

// GetForecastHandler handles GET requests for forecasts for a specific city
func (m *Manager) GetForecastHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	f, err := m.FM.Get(vars["id"])
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fc := &Forecast{
		CityID: f.CityID.Hex(),
		Max:    f.Max,
		Min:    f.Min,
		Sample: f.Sample,
	}

	b, err := json.Marshal(fc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
