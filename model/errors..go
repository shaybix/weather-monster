package model

import "errors"

var (
	// ErrNotFound describes an error where there are no document(s) in the database found
	ErrNotFound = errors.New("document(s) not found")
	// ErrAlreadyExists describes an error where a document already exists in the database
	ErrAlreadyExists = errors.New("document already exists")
)
