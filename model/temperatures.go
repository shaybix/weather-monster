package model

import (
	"database/sql"
	"time"
)

// Temperature describes a the temperature of any given day
type Temperature struct {
	ID        int64
	CityID    int64
	Min       int64
	Max       int64
	Timestamp int64
}

// NewTemperature describes a new temperature to be added for a city
type NewTemperature struct {
	CityID int64
	Min    int64
	Max    int64
}

// TemperatureManager describes a temperature model manager
type TemperatureManager struct {
	DB *sql.DB
}

// Create creates a temperature entry in the database
func (tm *TemperatureManager) Create(tf *NewTemperature) (*Temperature, error) {

	var temp Temperature

	sqlStmt := `
	INSERT INTO temperatures 
	(city_id, min, max, timestamp) 
	VALUES($1, $2, $3, $4) 
	RETURNING ID, min, max, timestamp, city_id;
	`
	if err := tm.DB.QueryRow(sqlStmt, tf.CityID, tf.Min, tf.Max, time.Now().Unix()).
		Scan(&temp.ID, &temp.Min, &temp.Max, &temp.Timestamp, &temp.CityID); err != nil {
		return nil, err
	}

	return &temp, nil
}

// NewTemperatureManager returns a new TemperatureManager
func NewTemperatureManager(db *sql.DB) *TemperatureManager {
	return &TemperatureManager{db}
}
