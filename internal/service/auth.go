package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/MorZLE/auth/internal/generate/jwtgen"
	"github.com/MorZLE/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user exists")
)

type UserSaver interface {
	SaveUser(ctx context.Context, login string, pswdHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvide interface {
	App(ctx context.Context, appID int32) (models.App, error)
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

	log := s.log.With(slog.String("op", op),
		slog.String("login", login))
	log.Info("login user")

	user, err := s.usrProvider.User(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.log.Warn("user not found", slog.String("login", login),
				slog.String("op", op),
				slog.String("err", err.Error()))
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("error get user %s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Error("invalid password", slog.String("err", err.Error()))

		return "", fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
	}

	app, err := s.appProvider.App(ctx, appID)
	if err != nil {
		s.log.Error("error get app", slog.String("err", err.Error()))
		return "", storage.ErrAppNotFound
	}

	token, err = jwtgen.NewJWT(user, app, s.tokenTTL)
	if err != nil {
		s.log.Error("error generate token", slog.String("err", err.Error()))
		return "", fmt.Errorf("error generate token %s: %w", op, err)
	}

	log.Info("user login success")

	return token, nil
}

func (s *Auth) RegisterNewUser(ctx context.Context, login string, password string) (userid int64, err error) {
	const op = "Auth.RegisterNewUser"

	log := s.log.With(slog.String("op", op), slog.String("login", login))

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed generate passhash")
		return 0, fmt.Errorf("failed generate passhash %s: %w", op, err)
	}

	uid, err := s.usrSaver.SaveUser(ctx, login, passhash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user exists", slog.String("err", err.Error()))
			return 0, ErrUserExists
		}
		return 0, fmt.Errorf("error save user %s: %w", op, err)
	}

	log.Info("register user")
	return uid, nil
}

func (s *Auth) CheckIsAdmin(ctx context.Context, userid int64) (bool, error) {
	const op = "auth.checkIsAdmin"

	log := s.log.With(slog.String("op", op), slog.Int64("userid", userid))

	ch, err := s.usrProvider.IsAdmin(ctx, userid)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return false, ErrInvalidCredentials
		}
		log.Error("error check is admin", slog.String("err", err.Error()))
		return false, err
	}
	log.Info("check is admin")

	return ch, nil
}
