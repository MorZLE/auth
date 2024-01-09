package cerror

import "github.com/pkg/errors"

var (
	ErrNotRights   = errors.New("not enough rights")
	ErrInternalErr = errors.New("internal err")
)

var (
	ErrAppNotFound        = errors.New("app not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user exists")
	ErrAppExists          = errors.New("app exists")
)
