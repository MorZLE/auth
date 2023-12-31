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

func (s *Storage) SaveUser(ctx context.Context, login string, pswdHash []byte, appid int32) (uid int64, err error) {
	const op = "sqlite.SaveUser"
	query := "INSERT INTO users (login, passHash,app_id) VALUES (?, ?, ?)"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, login, pswdHash, appid)
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

func (s *Storage) User(ctx context.Context, login string, appid int32) (models.User, error) {
	var user models.User
	const op = "sqlite.User"
	query := "SELECT id, login,passHash FROM users WHERE login = ? and app_id = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, login, appid)
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

func (s *Storage) IsAdmin(ctx context.Context, userID int32, appID int32) (models.Admin, error) {
	var res models.Admin
	const op = "sqlite.IsAdmin"
	query := "SELECT id,user_id,lvl,app_id FROM admins WHERE user_id = ? and app_id = ?"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return res, fmt.Errorf("%s: %w ", op, err)
	}
	row := stmt.QueryRowContext(ctx, userID, appID)

	err = row.Scan(&res.Id, &res.UserID, &res.Lvl, &res.AppID)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sql.ErrNoRows {
			return res, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return res, fmt.Errorf("%s: %w ", op, err)
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

	row := stmt.QueryRowContext(ctx, appID)
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

func (s *Storage) CreateAdmin(ctx context.Context, login string, lvl int32, appID int32) (uid int64, err error) {
	const op = "storage.CreateAdmin"
	query := "INSERT INTO admins (user_id, lvl, app_id) SELECT id, ?, ? FROM users WHERE login = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w ", op, err)
	}

	res, err := stmt.ExecContext(ctx, lvl, appID, login)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w ", op, err)
	}

	uid, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w ", op, err)
	}
	return uid, nil
}

func (s *Storage) DeleteAdmin(ctx context.Context, login string) (res bool, err error) {
	const op = "storage.DeleteAdmin"
	query := "delete * from admins where user_id = (select id from users where login= ?)"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w ", op, err)
	}
	_, err = stmt.ExecContext(ctx, login)
	if err != nil {
		return false, fmt.Errorf("%s: %w ", op, err)
	}
	return true, err
}
func (s *Storage) AddApp(ctx context.Context, name, secret string) (int32, error) {
	const op = "storage.AddApp"

	query := "INSERT INTO apps (name,secret) VALUES(?,?)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, name, secret)
	if err != nil {
		var errSql sqlite3.Error
		if errors.As(err, &errSql) && errSql.ExtendedCode == sql.ErrNoRows {
			return 0, storage.ErrAppExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	uid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(uid), nil
}
