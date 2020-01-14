package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shaybix/weather-monster/model"
	"net/http"
	"strconv"
)

// City describes a city and its' location in the world
type City struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Version   string  `json:"version"`
}

// CreateCityHandler handles a POST request to create a city
func (m *Manager) CreateCityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	city, err := m.CM.Create(&model.NewCity{
		Name:      r.FormValue("name"),
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		if err == model.ErrAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&City{
		ID:        city.ID.Hex(),
		Name:      city.Name,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Version:   city.Version.Hex(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

// UpdateCityHandler handles a PATCH request to update a city
func (m *Manager) UpdateCityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	city, err := m.CM.Update(&model.CityUpdate{
		ID:        vars["id"],
		Name:      r.FormValue("name"),
		Latitude:  latitude,
		Longitude: longitude,
		Version:   r.FormValue("version"),
	})
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&City{
		ID:        city.ID.Hex(),
		Name:      city.Name,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Version:   city.Version.Hex(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// DeleteCityHandler handles a DELETE request to delete a city
func (m *Manager) DeleteCityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	city, err := m.CM.Delete(vars["id"])
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&City{
		ID:        city.ID.Hex(),
		Name:      city.Name,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Version:   city.Version.Hex(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
