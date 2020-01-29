package service

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shaybix/weather-monster/model"
)

func withServiceTestDB(f func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T), t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	f(db, mock, t)
}

// Manager ...
type Manager struct {
	CM *model.CityManager
	FM *model.ForecastManager
	TM *model.TemperatureManager
	WM *model.WebhookManager
}

// NewServiceManager ...
func NewServiceManager(db *sql.DB) *Manager {
	return &Manager{
		CM: model.NewCityManager(db),
		FM: model.NewForecastManager(db),
		TM: model.NewTemperatureManager(db),
		WM: model.NewWebhookManager(db),
	}
}
