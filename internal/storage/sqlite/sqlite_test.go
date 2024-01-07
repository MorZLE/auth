package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/MorZLE/auth/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"reflect"
	"testing"
)

func TestStorage_AddApp(t *testing.T) {

	db, remove := goTestDB(sqlite)
	defer remove()

	checkres := func(t *testing.T, db *sql.DB, app models.App) bool {
		var gotRes models.App
		query := "SELECT id,name,secret FROM apps WHERE id = ?"

		stmt, err := db.Prepare(query)
		if err != nil {
			t.Errorf("AddApp() error = %v", err)
		}

		row := stmt.QueryRowContext(context.Background(), app.ID)
		err = row.Scan(&gotRes.ID, &gotRes.Name, &gotRes.Secret)
		if err != nil {
			t.Errorf("AddApp() error = %v", err)
		}

		if !reflect.DeepEqual(app, gotRes) {
			t.Errorf("AddApp() got = %v, want %v", gotRes, app)
		}
		return true
	}

	type args struct {
		name   string
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		check   bool
		wantErr error
	}{
		{
			name: "positive_1",
			args: args{
				name:   "awdawf",
				secret: "secret",
			},
			want:    int32(1),
			wantErr: nil,
			check:   true,
		},
		{
			name: "positive_2",
			args: args{
				name:   "morzle.com",
				secret: "sef23fresef",
			},
			want:    int32(2),
			wantErr: nil,
			check:   true,
		},
		{
			name: "negative_1",
			args: args{
				name:   "morzle.com",
				secret: "sef23awd",
			},
			want:    int32(0),
			wantErr: storage.ErrUniqueApp,
			check:   false,
		},
		{
			name: "negative_2",
			args: args{
				name:   "morzle.com",
				secret: "sef23awd",
			},
			want:    int32(0),
			wantErr: storage.ErrUniqueApp,
			check:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: db,
			}
			got, err := s.AddApp(context.Background(), tt.args.name, tt.args.secret)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddApp() got = %v, want %v", got, tt.want)
			}
			if tt.check {
				if !checkres(t, db, models.App{ID: int64(got), Name: tt.args.name, Secret: tt.args.secret}) {
					t.Errorf("AddApp() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestStorage_App(t *testing.T) {

	db, remove := goTestDB(sqlite)
	defer remove()

	apps := []models.App{
		{ID: 1, Name: "awdawf", Secret: "secret"},
		{ID: 2, Name: "morzGEAHle.com", Secret: "sef23fresef"},
		{ID: 3, Name: "morzASZDe.com", Secret: "awdtrdehdf"},
		{ID: 4, Name: "morEWFzle.com", Secret: "sef23aSEf432gawd"},
		{ID: 5, Name: "morzle.com", Secret: "sef2aerbvb45whnw3awd"},
	}

	checkres := func(t *testing.T, db *sql.DB, app models.App) {
		query := "INSERT INTO apps (name,secret) VALUES(?,?)"
		stmt, err := db.Prepare(query)
		if err != nil {
			t.Errorf("AddApp() error = %v", err)
		}
		_, err = stmt.ExecContext(context.Background(), app.Name, app.Secret)
		if err != nil {
			t.Errorf("AddApp() error = %v", err)
		}
	}

	for _, app := range apps {
		checkres(t, db, app)
	}

	func(t *testing.T) {
		query := "delete from apps where id = ?"
		stmt, err := db.Prepare(query)
		if err != nil {
			t.Errorf("AddApp() error = %v", err)
		}
		_, err = stmt.ExecContext(context.Background(), 3)
		if err != nil {
			t.Errorf("AddApp() error = %v", err)
		}
	}(t)

	type args struct {
		appID int32
	}
	tests := []struct {
		name    string
		args    args
		want    models.App
		wantErr error
	}{
		{
			name: "positive_1",
			args: args{
				appID: 1,
			},
			want:    apps[0],
			wantErr: nil,
		},
		{
			name: "positive_2",
			args: args{
				appID: 4,
			},
			want:    apps[3],
			wantErr: nil,
		},
		{
			name: "positive_3",
			args: args{
				appID: 5,
			},
			want:    apps[4],
			wantErr: nil,
		},
		{
			name: "negative_1",
			args: args{
				appID: 6,
			},
			want:    models.App{},
			wantErr: storage.ErrAppNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: db,
			}
			got, err := s.App(context.Background(), tt.args.appID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("AddApp() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestStorage_CreateAdmin(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		login string
		lvl   int32
		appID int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantUid int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: tt.fields.db,
			}
			gotUid, err := s.CreateAdmin(tt.args.ctx, tt.args.login, tt.args.lvl, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUid != tt.wantUid {
				t.Errorf("CreateAdmin() gotUid = %v, want %v", gotUid, tt.wantUid)
			}
		})
	}
}

func TestStorage_DeleteAdmin(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: tt.fields.db,
			}
			gotRes, err := s.DeleteAdmin(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("DeleteAdmin() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestStorage_IsAdmin(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int32
		appID  int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Admin
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: tt.fields.db,
			}
			got, err := s.IsAdmin(tt.args.ctx, tt.args.userID, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsAdmin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_SaveUser(t *testing.T) {

	db, remove := goTestDB(sqlite)
	defer remove()

	checkDB := func(login string, appid int32) error {
		var resUser models.User
		query := "SELECT id, login,passHash FROM users WHERE login = ? and app_id = ?"

		stmt, err := db.Prepare(query)
		if err != nil {
			return err
		}
		res := stmt.QueryRowContext(context.Background(), login, appid)
		err = res.Scan(&resUser.ID, &resUser.Login, &resUser.PassHash)
		if err != nil {
			return err
		}
		return nil
	}

	type args struct {
		login    string
		pswdHash []byte
		appid    int32
	}
	tests := []struct {
		name    string
		args    args
		wantUid int64
		wantErr bool
		error   error
	}{
		{
			name: "positive_1",
			args: args{
				login:    "awd",
				pswdHash: []byte("awd"),
				appid:    1,
			},
			wantUid: 1,
			wantErr: false,
			error:   nil,
		},
		{
			name: "positive_2",
			args: args{
				login:    "awfawfawf",
				pswdHash: []byte("awawfhgrthd"),
				appid:    2,
			},
			wantUid: 2,
			wantErr: false,
			error:   nil,
		},
		{
			name: "negative_1",
			args: args{
				login:    "awfawfawf",
				pswdHash: []byte("awawfhgrthd"),
				appid:    2,
			},
			wantUid: 0,
			wantErr: true,
			error:   storage.ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: db,
			}
			gotUid, err := s.SaveUser(context.Background(), tt.args.login, tt.args.pswdHash, tt.args.appid)
			if !errors.Is(err, tt.error) {
				t.Errorf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotUid != tt.wantUid {
				t.Errorf("SaveUser() gotUid = %v, want %v", gotUid, tt.wantUid)
			}
			err = checkDB(tt.args.login, tt.args.appid)
			if err != nil {
				t.Errorf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestStorage_User(t *testing.T) {

	db, closeDB := goTestDB(sqlite)
	defer closeDB()

	createUser := func(db *sql.DB, user models.User, appID int32) {
		query := "INSERT INTO users (login, passHash,app_id) VALUES (?, ?, ?)"

		stmt, err := db.Prepare(query)
		if err != nil {
			panic(err)
		}
		_, err = stmt.ExecContext(context.Background(), user.Login, user.PassHash, appID)
	}

	users := []models.User{
		{ID: 1, Login: "test", PassHash: []byte("123")},
		{ID: 2, Login: "awd", PassHash: []byte("125323")},
		{ID: 3, Login: "tedhe5hst", PassHash: []byte("122343")},
		{ID: 4, Login: "tezdbe4gst", PassHash: []byte("122345783")},
	}
	appIDS := []int32{1, 2, 3, 4}

	for i, user := range users {
		createUser(db, user, appIDS[i])
	}

	type args struct {
		login string
		appid int32
	}
	tests := []struct {
		name    string
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "positive_1",
			args: args{
				login: "test",
				appid: 1,
			},
			want:    users[0],
			wantErr: false,
		},
		{
			name: "positive_2",
			args: args{
				login: "tedhe5hst",
				appid: 3,
			},
			want:    users[2],
			wantErr: false,
		},
		{
			name: "negative_1",
			args: args{
				login: "tedhe5hst",
				appid: 1,
			},
			want:    models.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: db,
			}
			got, err := s.User(context.Background(), tt.args.login, tt.args.appid)
			if (err != nil) != tt.wantErr {
				t.Errorf("User() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User() got = %v, want %v", got, tt.want)
			}
		})
	}
}

const sqlite = "sqlite3"

func goTestDB(vendor string) (*sql.DB, func()) {
	switch vendor {
	case sqlite:
		return createTestSqliteDB("D:/Golang/auth/storage/test.db", "D:/Golang/auth/migrations/", "")
	}
	panic("unknown vendor db")
	return nil, nil
}

func createTestSqliteDB(storagePath, migratorPath, migratorTable string) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		panic(err)
	}

	m, err := migrate.New(
		"file:"+migratorPath,
		fmt.Sprintf("sqlite3://%s?x-migration-table=%s", storagePath, migratorTable),
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}
	err, _ = m.Close()
	if err != nil {
		panic(err)
	}

	return db, func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}

		err = os.Remove(storagePath)
		if err != nil {
			panic(err)
		}
	}
}
