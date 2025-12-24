package repository

import "errors"

var (
	// ErrNotFound is returned when record does not exist.
	ErrNotFound = errors.New("not found")
	// ErrNotEnoughStock is returned when there is not enough product quantity.
	ErrNotEnoughStock = errors.New("not enough stock")
)
