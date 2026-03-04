package infrastructure

import "errors"

var (
	ErrItemNotFound = errors.New("Item not found")
	ErrBadArgument  = errors.New("Wrong argument")
)
