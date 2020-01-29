package model

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

// Webhook describes a webhook for subscribing to a city's temperatures
type Webhook struct {
	ID          int64
	CityID      int64
	CallbackURL string
}

// NewWebhook describes a new webhook to be created
type NewWebhook struct {
	CityID      int64
	CallbackURL string
}

// WebhookManager describes a webhook model manager
type WebhookManager struct {
	db *sql.DB
}

// Create creates a new webhook for a given city
func (w *WebhookManager) Create(nw *NewWebhook) (*Webhook, error) {

	var wh Webhook
	sqlStmt := `
	INSERT INTO webhooks 
	(city_id, callback_url) 
	VALUES($1, $2)
	RETURNING ID, city_id, callback_url;`
	if err := w.db.QueryRow(sqlStmt, nw.CityID, nw.CallbackURL).
		Scan(&wh.ID, &wh.CityID, &wh.CallbackURL); err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	return &wh, nil
}

// Get gets all webhooks associated with a city
func (w *WebhookManager) Get(cityID int64) ([]*Webhook, error) {

	var webhooks []*Webhook
	sqlStmt := `
	SELECT * FROM webhooks 
	WHERE city_id = $1
	RETURNING ID, city_id, callback_url;`

	rows, err := w.db.Query(sqlStmt, cityID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "20000" {
				return nil, ErrNotFound
			}
		}
		return nil, err
	}

	for rows.Next() {
		var wh Webhook
		if err := rows.Scan(&wh.ID, &wh.CityID, &wh.CallbackURL); err != nil {
			log.Println(err)
			continue
		}

		webhooks = append(webhooks, &wh)
	}

	if len(webhooks) == 0 {
		return nil, ErrNotFound
	}

	return webhooks, nil
}

// Delete deletes a webhook
func (w *WebhookManager) Delete(id int64) (*Webhook, error) {

	var wh Webhook
	sqlStmt := `
	DELETE FROM webhooks 
	WHERE ID = $1
	RETURNING ID, city_id, callback_url;`
	if err := w.db.QueryRow(sqlStmt, id).
		Scan(&wh.ID, &wh.CityID, &wh.CallbackURL); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "20000" {
				return nil, ErrNotFound
			}
		}
		return nil, err
	}

	return &wh, nil
}

// NewWebhookManager returns a new WebhookManager
func NewWebhookManager(db *sql.DB) *WebhookManager {
	wm := &WebhookManager{db}
	return wm
}
