package service

import (
	"context"
	"github.com/MorZLE/auth/internal/domain/models"
	"log/slog"
	"time"
)

type UserSaver interface {
	SaveUser(ctx context.Context, login string, pswdHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvide interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

// NewAuth возвращает новый экземпляр сервиса
func NewAuth(log *slog.Logger,
	usrProvider UserProvider,
	usrSaver UserSaver,
	appProvider AppProvide,
	tokenTTL time.Duration) *Auth {

	return &Auth{log: log, usrProvider: usrProvider, usrSaver: usrSaver, appProvider: appProvider, tokenTTL: tokenTTL}
}

type Auth struct {
	log         *slog.Logger
	usrProvider UserProvider
	usrSaver    UserSaver
	appProvider AppProvide
	tokenTTL    time.Duration
}

func (s *Auth) LoginUser(ctx context.Context, login string, password string, appID int32) (token string, err error) {
	const op = "Auth.LoginUser"

	log := s.log.With(slog.String("op", op), slog.String("login", login))

	log.Info("login user")
	return "", nil
}

func (s *Auth) RegisterNewUser(ctx context.Context, login string, password string) (userid int64, err error) {
	const op = "Auth.RegisterNewUser"

	log := s.log.With(slog.String("op", op), slog.String("login", login))

	log.Info("register user")

	return 0, nil
}

func (s *Auth) CheckIsAdmin(userid int32) (bool, error) {
	return false, nil
}
