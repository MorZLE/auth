package controller

import (
	"context"
	"github.com/MorZLE/auth/internal/domain/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Auth
type Auth interface {
	LoginUser(ctx context.Context, login string, password string, appID int32) (token string, err error)
	RegisterNewUser(ctx context.Context, login string, password string, appid int32) (userid int64, err error)
	CheckIsAdmin(ctx context.Context, userid int32, appID int32) (models.Admin, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=AuthAdmin
type AuthAdmin interface {
	CreateAdmin(ctx context.Context, login string, lvl int32, key string, appid int32) (userid int64, err error)
	DeleteAdmin(ctx context.Context, login string, key string) (res bool, err error)
	AddApp(ctx context.Context, name, secret, key string) (userid int32, err error)
}
