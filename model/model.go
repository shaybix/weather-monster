package model

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func withTestDB(f func(db *sql.DB, mock sqlmock.Sqlmock, t *testing.T), t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	f(db, mock, t)
}
