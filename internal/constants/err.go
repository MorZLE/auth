package constants

import "github.com/pkg/errors"

var (
	ErrNotRights   = errors.New("not enough rights")
	ErrInternalErr = errors.New("internal err")
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user exists")
	ErrAppExists          = errors.New("app exists")
)
