package repository

import "errors"

var (
	// ErrRecordNotFound record not found
	ErrRecordNotFound = errors.New("record not found")
	// ErrNoRecordUpdated no record updated
	ErrNoRecordUpdated = errors.New("no record updated")
	// ErrNoRecordDeleted no record deleted
	ErrNoRecordDeleted = errors.New("no record deleted")
)

