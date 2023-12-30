package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/MorZLE/auth/internal/constants"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/MorZLE/auth/internal/generate/jwtgen"
	"github.com/MorZLE/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserSaver interface {
	SaveUser(ctx context.Context, login string, pswdHash []byte, appid int32) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	IsAdmin(ctx context.Context, userID int32) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

type AdminProvider interface {
	CreateAdmin(ctx context.Context, login string, lvl int32) (uid int64, err error)
	DeleteAdmin(ctx context.Context, login string) (uid int64, err error)
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

	user, err := s.usrProvider.User(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.log.Warn("user not found", slog.String("login", login),
				slog.String("op", op),
				slog.String("err", err.Error()))
			return "", constants.ErrInvalidCredentials
		}
		return "", fmt.Errorf("error get user %s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Error("invalid password", slog.String("err", err.Error()))

		return "", fmt.Errorf("%s : %w", op, constants.ErrInvalidCredentials)
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
			return 0, constants.ErrUserExists
		}
		return 0, fmt.Errorf("error save user %s: %w", op, err)
	}

	log.Info("register user")
	return uid, nil
}

func (s *Auth) CheckIsAdmin(ctx context.Context, userid int32) (bool, error) {
	const op = "auth.checkIsAdmin"

	log := s.log.With(slog.String("op", op), slog.Int64("userid", int64(userid)))

	ch, err := s.usrProvider.IsAdmin(ctx, userid)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return false, constants.ErrInvalidCredentials
		}
		log.Error("error check is admin", slog.String("err", err.Error()))
		return false, err
	}
	log.Info("check is admin")

	return ch, nil
}

func (s *Auth) CreateAdmin(ctx context.Context, login string, lvl int32, key string) (userid int64, err error) {
	const op = "auth.CreateAdmin"
	log := s.log.With(slog.String("op", op), slog.String("login", login), slog.Int("lvl", int(lvl)))
	if !checkKeyAdmin(key) {
		return 0, constants.ErrNotRights
	}

	uid, err := s.admProvider.CreateAdmin(ctx, login, lvl)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return 0, constants.ErrInvalidCredentials
		}
		log.Error("error createAdmin is admin", slog.String("err", err.Error()))
		return 0, err
	}

	log.Info(fmt.Sprintf("change root admin %s", login))
	return uid, nil
}

func (s *Auth) DeleteAdmin(ctx context.Context, login string, key string) (userid int64, err error) {
	const op = "auth.DeleteAdmin"

	if !checkKeyAdmin(key) {
		return 0, constants.ErrNotRights
	}

	log := s.log.With(slog.String("op", op), slog.String("login", login))

	uid, err := s.admProvider.DeleteAdmin(ctx, login)
	if err != nil {
		log.Error("error DeleteAdmin", slog.String("err", err.Error()))
		if errors.Is(err, storage.ErrUserNotFound) {
			return 0, constants.ErrInvalidCredentials
		}
		return 0, constants.ErrInternalErr
	}

	log.Info(fmt.Sprintf("delete admin %s", login))

	return uid, nil
}
func (s *Auth) AddApp(ctx context.Context, name, secret, key string) (userid int32, err error) {
	const op = "auth.AddApp"

	if !checkKeyAdmin(key) {
		return 0, constants.ErrNotRights
	}

	log := s.log.With(slog.String("op", op), slog.String("name", name))

	uid, err := s.admProvider.AddApp(ctx, name, secret)
	if err != nil {
		log.Error("error AddApp", slog.String("err", err.Error()))
		if errors.Is(err, storage.ErrAppExists) {
			return 0, constants.ErrUserExists
		}
		return 0, constants.ErrInternalErr
	}
	log.Info(fmt.Sprintf("add app: %s", name))

	return uid, nil
}

func checkKeyAdmin(key string) bool {
	return key == "pppp"
}
