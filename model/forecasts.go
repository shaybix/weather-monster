package model

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

// Forecast describes a the forecast of a city with the average minimum and maximum temperature in the past 24h
type Forecast struct {
	CityID int64
	Min    int64
	Max    int64
	Sample int64
}

// ForecastManager describes a forecast model manager
type ForecastManager struct {
	DB *sql.DB
}

// Get returns the forecast of a city
func (fm *ForecastManager) Get(cid int64) (*Forecast, error) {

	sqlStmt := `
	SELECT min, max FROM temperatures 
	WHERE city_id = $1
	`
	rows, err := fm.DB.Query(sqlStmt, cid)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "20000" {
				return nil, ErrNotFound
			}
		}
		return nil, err
	}
	defer rows.Close()

	var sample int64
	var mins []int64
	var maxs []int64

	forecast := &Forecast{
		CityID: cid,
	}
	for rows.Next() {
		var temp Temperature
		if err := rows.Scan(&temp.Min, &temp.Max); err != nil {
			log.Println(err)
			continue
		}

		mins = append(mins, temp.Min)
		maxs = append(maxs, temp.Max)

		sample++

	}

	if sample == 0 {
		return forecast, nil
	}

	forecast.Sample = sample
	forecast.Min = sum(mins) / int64(len(mins))
	forecast.Min = sum(maxs) / int64(len(maxs))

	return forecast, nil
}

func sum(temps []int64) int64 {
	var total int64
	for _, temp := range temps {
		total = total + temp
	}

	return total
}

// NewForecastManager returns a new ForecastManager
func NewForecastManager(db *sql.DB) *ForecastManager {
	return &ForecastManager{
		DB: db,
	}
}
