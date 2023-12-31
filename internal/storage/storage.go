package storage

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user exists")
	ErrAppExists    = errors.New("app exist")
	ErrAppNotFound  = errors.New("app not found")
)
