package infrastructure

import "errors"

var (
	ErrItemNotFound = errors.New("item not found")
	ErrIDNotFound   = errors.New("id not found")
)
