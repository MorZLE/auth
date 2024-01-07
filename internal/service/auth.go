package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/MorZLE/auth/internal/domain/cerror"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/MorZLE/auth/internal/generate/jwtgen"
	"github.com/MorZLE/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=UserSaver
type UserSaver interface {
	SaveUser(ctx context.Context, login string, pswdHash []byte, appid int32) (uid int64, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=UserProvider
type UserProvider interface {
	User(ctx context.Context, login string, appid int32) (models.User, error)
	IsAdmin(ctx context.Context, userID int32, appid int32) (models.Admin, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=AppProvider
type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=AdminProvider
type AdminProvider interface {
	CreateAdmin(ctx context.Context, login string, lvl int32, appID int32) (uid int64, err error)
	DeleteAdmin(ctx context.Context, login string) (res bool, err error)
	AddApp(ctx context.Context, name, secret string) (uid int32, err error)
}

// NewAuth возвращает новый экземпляр сервиса
func NewAuth(log *slog.Logger,
	usrProvider UserProvider,
	usrSaver UserSaver,
	appProvider AppProvider,
	admProvider AdminProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{log: log, usrProvider: usrProvider, usrSaver: usrSaver, appProvider: appProvider, admProvider: admProvider, tokenTTL: tokenTTL}
}

type Auth struct {
	log         *slog.Logger
	usrProvider UserProvider
	usrSaver    UserSaver
	appProvider AppProvider
	admProvider AdminProvider
	tokenTTL    time.Duration
}

func (s *Auth) LoginUser(ctx context.Context, login string, password string, appID int32) (token string, err error) {
	const op = "Auth.LoginUser"

	log := s.log.With(slog.String("op", op),
		slog.String("login", login))
	log.Info("login user")

	user, err := s.usrProvider.User(ctx, login, appID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.log.Warn("user not found", slog.String("login", login),
				slog.String("op", op),
				slog.String("err", err.Error()))
			return "", cerror.ErrInvalidCredentials
		}
		return "", fmt.Errorf("cerror get user %s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Error("invalid password", slog.String("err", err.Error()))

		return "", fmt.Errorf("%s : %w", op, cerror.ErrInvalidCredentials)
	}

	app, err := s.appProvider.App(ctx, appID)
	if err != nil {
		s.log.Error("cerror get app", slog.String("err", err.Error()))
		return "", storage.ErrAppNotFound
	}

	token, err = jwtgen.NewJWT(user, app, s.tokenTTL)
	if err != nil {
		s.log.Error("cerror generate token", slog.String("err", err.Error()))
		return "", fmt.Errorf("cerror generate token %s: %w", op, err)
	}

	log.Info("user login success")

	return token, nil
}

func (s *Auth) RegisterNewUser(ctx context.Context, login string, password string, appid int32) (userid int64, err error) {
	const op = "Auth.RegisterNewUser"

	log := s.log.With(slog.String("op", op), slog.String("login", login))

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed generate passhash")
		return 0, fmt.Errorf("failed generate passhash %s: %w", op, err)
	}

	uid, err := s.usrSaver.SaveUser(ctx, login, passhash, appid)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user exists", slog.String("err", err.Error()))
			return 0, cerror.ErrUserExists
		}
		return 0, fmt.Errorf("cerror save user %s: %w", op, err)
	}

	log.Info("register user")
	return uid, nil
}

func (s *Auth) CheckIsAdmin(ctx context.Context, userid int32, appid int32) (models.Admin, error) {
	const op = "auth.checkIsAdmin"

	log := s.log.With(slog.String("op", op), slog.Int64("userid", int64(userid)))

	res, err := s.usrProvider.IsAdmin(ctx, userid, appid)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return res, cerror.ErrInvalidCredentials
		}
		log.Error("cerror check is admin", slog.String("err", err.Error()))
		return res, cerror.ErrInternalErr
	}
	log.Info("check is admin")

	return res, nil
}

func (s *Auth) CreateAdmin(ctx context.Context, login string, lvl int32, key string, appID int32) (userid int64, err error) {
	const op = "auth.CreateAdmin"
	log := s.log.With(slog.String("op", op), slog.String("login", login), slog.Int("lvl", int(lvl)))
	if !checkKeyAdmin(key) {
		return 0, cerror.ErrNotRights
	}

	uid, err := s.admProvider.CreateAdmin(ctx, login, lvl, appID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return 0, cerror.ErrInvalidCredentials
		}
		log.Error("cerror createAdmin is admin", slog.String("err", err.Error()))
		return 0, cerror.ErrInternalErr
	}

	log.Info(fmt.Sprintf("change root admin %s", login))
	return uid, nil
}

func (s *Auth) DeleteAdmin(ctx context.Context, login string, key string) (res bool, err error) {
	const op = "auth.DeleteAdmin"

	if !checkKeyAdmin(key) {
		return false, cerror.ErrNotRights
	}

	log := s.log.With(slog.String("op", op), slog.String("login", login))

	uid, err := s.admProvider.DeleteAdmin(ctx, login)
	if err != nil {
		log.Error("cerror DeleteAdmin", slog.String("err", err.Error()))

		if errors.Is(err, storage.ErrUserNotFound) {
			return false, cerror.ErrInvalidCredentials
		}
		return false, cerror.ErrInternalErr
	}

	log.Info(fmt.Sprintf("delete admin %s", login))

	return uid, nil
}
func (s *Auth) AddApp(ctx context.Context, name, secret, key string) (userid int32, err error) {
	const op = "auth.AddApp"

	if !checkKeyAdmin(key) {
		return 0, cerror.ErrNotRights
	}

	log := s.log.With(slog.String("op", op), slog.String("name", name))

	uid, err := s.admProvider.AddApp(ctx, name, secret)
	if err != nil {
		log.Error("cerror AddApp", slog.String("err", err.Error()))
		if errors.Is(err, storage.ErrAppExists) {
			return 0, cerror.ErrUserExists
		}
		return 0, cerror.ErrInternalErr
	}
	log.Info(fmt.Sprintf("add app: %s", name))

	return uid, nil
}

func checkKeyAdmin(key string) bool {
	return key != ""
}
