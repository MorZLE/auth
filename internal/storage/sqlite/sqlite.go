package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/MorZLE/auth/internal/storage"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func NewStorage(dbPath string) (*Storage, error) {
	const op = "sqlite.NewStorage"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

type Storage struct {
	db *sql.DB
}

type UserSaver struct{}
type UserProvider struct{}
type AppProvide struct{}

func (s *Storage) SaveUser(ctx context.Context, login string, pswdHash []byte) (uid int64, err error) {
	const op = "sqlite.SaveUser"
	query := "INSERT INTO users (login, password) VALUES (?, ?)"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, login, pswdHash)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, login string) (models.User, error) {
	var user models.User
	const op = "sqlite.User"
	query := "SELECT id, login,passHash FROM users WHERE login = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, login)
	err = res.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sql.ErrNoRows {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int32) (bool, error) {
	var res bool
	const op = "sqlite.IsAdmin"
	query := "SELECT isadmin FROM users WHERE id = ?"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w ", op, err)
	}
	row := stmt.QueryRowContext(ctx, userID)

	err = row.Scan(&res)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sql.ErrNoRows {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w ", op, err)
	}

	return res, nil
}

func (s *Storage) App(ctx context.Context, appID int32) (models.App, error) {
	const op = "sqlite.App"
	var res models.App
	query := "SELECT id,name,secret FROM apps WHERE id = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return res, fmt.Errorf("%s: %w ", op, err)
	}

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(&res.ID, &res.Name, &res.Secret)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sql.ErrNoRows {
			return res, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return res, fmt.Errorf("%s: %w ", op, err)
	}
	return res, nil
}
