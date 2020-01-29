package model

import (
	"database/sql"

	"github.com/gobuffalo/uuid"
	"github.com/lib/pq"
)

// City describes a city in the world, e.g. Berlin
type City struct {
	ID        int64
	Name      string
	Latitude  float64
	Longitude float64
	Version   string
}

// NewCity describes the form values of a new city
type NewCity struct {
	Name      string
	Latitude  float64
	Longitude float64
}

// CityUpdate describes the form values of a city to be updated
type CityUpdate struct {
	ID        int64
	Name      string
	Latitude  float64
	Longitude float64
	Version   string
}

// CityManager describes a city model mannager
type CityManager struct {
	db *sql.DB
}

// Create creates a new non-existing entry of a city in the database
func (cm *CityManager) Create(nc *NewCity) (*City, error) {
	version, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	INSERT INTO cities
	(name, latitude, longitude, version) 
	VALUES($1,$2,$3,$4)
	RETURNING ID, name, latitude, longitude, version;
	`

	var city City
	if err := cm.db.QueryRow(sqlStmt, nc.Name, nc.Latitude, nc.Longitude, version.String()).
		Scan(&city.ID, &city.Name, &city.Latitude, &city.Longitude, &city.Version); err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	return &city, nil
}

// Update updates an existing city  in the database
func (cm *CityManager) Update(cu *CityUpdate) (*City, error) {
	newVersion, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	UPDATE cities
	SET name = $1, latitude = $2, longitude = $3, version = $4
	WHERE ID = $5 AND version = $6
	RETURNING ID, name, latitude, longitude, version;
	`

	var city City
	if err := cm.db.QueryRow(sqlStmt, cu.Name, cu.Latitude, cu.Longitude, newVersion.String(), cu.ID, cu.Version).
		Scan(&city.ID, &city.Name, &city.Latitude, &city.Longitude, &city.Version); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "20000" {
				return nil, ErrNotFound
			}
		}
		return nil, err
	}

	return &city, nil
}

// Delete deletes an existing city in the database
func (cm *CityManager) Delete(id int64) (*City, error) {
	sqlStmt := `
	DELETE FROM cities
	WHERE ID = $1
	RETURNING ID, name, latitude, longitude, version;
	`

	var city City
	if err := cm.db.QueryRow(sqlStmt, id).
		Scan(&city.ID, &city.Name, &city.Latitude, &city.Longitude, &city.Version); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "20000" {
				return nil, ErrNotFound
			}
		}
		return nil, err
	}

	return &city, nil
}

// NewCityManager returns a new CityManager
func NewCityManager(db *sql.DB) *CityManager {
	cm := &CityManager{db}

	return cm
}
