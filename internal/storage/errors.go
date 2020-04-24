package storage

import "errors"

var (
	ErrDuplicate = errors.New("already exists")
	ErrNotFound  = errors.New("not found error")
)
