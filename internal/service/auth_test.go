package service

import (
	"context"
	"errors"
	"github.com/MorZLE/auth/internal/domain/constants"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/MorZLE/auth/internal/service/mocks"
	"github.com/MorZLE/auth/internal/storage"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"reflect"
	"testing"
)

const keyAdmin = "test"

func TestAuth_AddApp(t *testing.T) {

	type mck func(s *mocks.AdminProvider)

	type args struct {
		name   string
		secret string
		key    string
	}
	tests := []struct {
		name       string
		args       args
		mck        mck
		wantUserid int32
		wantErr    error
	}{
		{
			name: "positive_1",
			args: args{
				name:   "test",
				secret: "test",
				key:    keyAdmin,
			},
			mck: func(s *mocks.AdminProvider) {
				s.On("AddApp", mock.Anything, "test", "test").Return(int32(1), nil)
			},
			wantUserid: int32(1),
			wantErr:    nil,
		},
		{
			name: "positive_2",
			args: args{
				name:   "qwreqwrqwr",
				secret: "qwqwr",
				key:    keyAdmin,
			},
			mck: func(s *mocks.AdminProvider) {
				s.On("AddApp", mock.Anything, "qwreqwrqwr", "qwqwr").Return(int32(1), nil)
			},
			wantUserid: int32(1),
			wantErr:    nil,
		},
		{
			name: "negative_1",
			args: args{
				name:   "qwreqwrqwr",
				secret: "teqwrst",
				key:    "",
			},
			mck:        func(s *mocks.AdminProvider) {},
			wantUserid: 0,
			wantErr:    constants.ErrNotRights,
		},
		{
			name: "negative_2",
			args: args{
				name:   "qwreqwrqwr",
				secret: "teqwrst",
				key:    keyAdmin,
			},
			mck: func(s *mocks.AdminProvider) {
				s.On("AddApp", mock.Anything, "qwreqwrqwr", "teqwrst").Return(int32(0), storage.ErrAppExists)
			},
			wantUserid: 0,
			wantErr:    constants.ErrUserExists,
		},
		{
			name: "negative_3",
			args: args{
				name:   "qwreqwrqwr",
				secret: "teqwrst",
				key:    keyAdmin,
			},
			mck: func(s *mocks.AdminProvider) {
				s.On("AddApp", mock.Anything, "qwreqwrqwr", "teqwrst").Return(int32(0), errors.ErrUnsupported)
			},
			wantUserid: 0,
			wantErr:    constants.ErrInternalErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := mocks.NewAdminProvider(t)
			tt.mck(sqlite)

			s := &Auth{
				log:         slog.With(slog.String("service", "auth")),
				admProvider: sqlite,
			}
			gotUserid, err := s.AddApp(context.Background(), tt.args.name, tt.args.secret, tt.args.key)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserid != tt.wantUserid {
				t.Errorf("AddApp() gotUserid = %v, want %v", gotUserid, tt.wantUserid)
			}
		})
	}
}

func TestAuth_CheckIsAdmin(t *testing.T) {
	type mck func(s *mocks.UserProvider)

	type args struct {
		userid int32
		appid  int32
	}
	tests := []struct {
		name    string
		args    args
		mck     mck
		want    models.Admin
		wantErr error
	}{
		{
			name: "positive_1",
			args: args{
				userid: int32(2),
				appid:  int32(1),
			},
			mck: func(s *mocks.UserProvider) {
				s.On("IsAdmin", mock.Anything, int32(2), int32(1)).
					Return(models.Admin{Lvl: int32(2)}, nil)
			},
			want: models.Admin{
				Lvl: int32(2),
			},
			wantErr: nil,
		},
		{
			name: "positive_2",
			args: args{
				userid: int32(1324546134),
				appid:  int32(56),
			},
			mck: func(s *mocks.UserProvider) {
				s.On("IsAdmin", mock.Anything, int32(1324546134), int32(56)).
					Return(models.Admin{Lvl: int32(2)}, nil)
			},
			want: models.Admin{
				Lvl: int32(2),
			},
			wantErr: nil,
		},
		{
			name: "negative_1",
			args: args{
				userid: int32(1324546134),
				appid:  int32(56),
			},
			mck: func(s *mocks.UserProvider) {
				s.On("IsAdmin", mock.Anything, int32(1324546134), int32(56)).
					Return(models.Admin{}, storage.ErrUserNotFound)
			},
			want:    models.Admin{},
			wantErr: constants.ErrInvalidCredentials,
		},
		{
			name: "negative_2",
			args: args{
				userid: int32(1324546134),
				appid:  int32(56),
			},
			mck: func(s *mocks.UserProvider) {
				s.On("IsAdmin", mock.Anything, int32(1324546134), int32(56)).
					Return(models.Admin{}, errors.ErrUnsupported)
			},
			want:    models.Admin{},
			wantErr: constants.ErrInternalErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := mocks.NewUserProvider(t)
			tt.mck(sqlite)
			s := &Auth{
				log:         slog.With(slog.String("service", "auth")),
				usrProvider: sqlite,
			}
			got, err := s.CheckIsAdmin(context.Background(), tt.args.userid, tt.args.appid)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CheckIsAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckIsAdmin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_CreateAdmin(t *testing.T) {
	type mck func(m *mocks.AdminProvider)

	type args struct {
		login string
		lvl   int32
		key   string
		appID int32
	}
	tests := []struct {
		name       string
		mck        mck
		args       args
		wantUserid int64
		wantErr    error
	}{
		{
			name: "positive_1",
			mck: func(m *mocks.AdminProvider) {
				m.On("CreateAdmin", mock.Anything, "sefsef", int32(2), int32(1)).Return(int64(1), nil)
			},
			args: args{
				login: "sefsef",
				lvl:   int32(2),
				key:   keyAdmin,
				appID: int32(1),
			},
			wantUserid: int64(1),
			wantErr:    nil,
		},
		{
			name: "positive_2",
			mck: func(m *mocks.AdminProvider) {
				m.On("CreateAdmin", mock.Anything, "awfawfawfawf", int32(2342), int32(143)).Return(int64(234652346), nil)
			},
			args: args{
				login: "awfawfawfawf",
				lvl:   int32(2342),
				key:   keyAdmin,
				appID: int32(143),
			},
			wantUserid: int64(234652346),
			wantErr:    nil,
		},
		{
			name: "negative_1",
			mck: func(m *mocks.AdminProvider) {
				m.On("CreateAdmin", mock.Anything, "awfawfawfawf", int32(2342), int32(143)).Return(int64(0), storage.ErrUserNotFound)
			},
			args: args{
				login: "awfawfawfawf",
				lvl:   int32(2342),
				key:   keyAdmin,
				appID: int32(143),
			},
			wantUserid: int64(0),
			wantErr:    constants.ErrInvalidCredentials,
		},
		{
			name: "negative_2",
			mck: func(m *mocks.AdminProvider) {
				m.On("CreateAdmin", mock.Anything, "awfawfawfawf", int32(2342), int32(143)).Return(int64(0), errors.ErrUnsupported)
			},
			args: args{
				login: "awfawfawfawf",
				lvl:   int32(2342),
				key:   keyAdmin,
				appID: int32(143),
			},
			wantUserid: int64(0),
			wantErr:    constants.ErrInternalErr,
		},
		{
			name: "invalid_key",
			mck:  func(m *mocks.AdminProvider) {},
			args: args{
				login: "awfawfawfawf",
				lvl:   int32(2342),
				key:   "",
				appID: int32(143),
			},
			wantUserid: int64(0),
			wantErr:    constants.ErrNotRights,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := mocks.NewAdminProvider(t)
			tt.mck(sqlite)
			s := &Auth{
				log:         slog.With(slog.String("service", "auth")),
				admProvider: sqlite,
			}
			gotUserid, err := s.CreateAdmin(context.Background(), tt.args.login, tt.args.lvl, tt.args.key, tt.args.appID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserid != tt.wantUserid {
				t.Errorf("CreateAdmin() gotUserid = %v, want %v", gotUserid, tt.wantUserid)
			}
		})
	}
}

func TestAuth_DeleteAdmin(t *testing.T) {
	type mck func(m *mocks.AdminProvider)

	type args struct {
		login string
		key   string
	}

	tests := []struct {
		name    string
		args    args
		mck     mck
		wantRes bool
		wantErr error
	}{
		{
			name: "positive_1",
			mck: func(m *mocks.AdminProvider) {
				m.On("DeleteAdmin", mock.Anything, "sefsef").Return(true, nil)
			},
			args: args{
				login: "sefsef",
				key:   keyAdmin,
			},
			wantRes: true,
			wantErr: nil,
		},
		{
			name: "positive_2",
			mck: func(m *mocks.AdminProvider) {
				m.On("DeleteAdmin", mock.Anything, "awdgresbh").Return(true, nil)
			},
			args: args{
				login: "awdgresbh",
				key:   keyAdmin,
			},
			wantRes: true,
			wantErr: nil,
		},
		{
			name: "invalid_key",
			mck:  func(m *mocks.AdminProvider) {},
			args: args{
				login: "awdgresbh",
				key:   "",
			},
			wantRes: false,
			wantErr: constants.ErrNotRights,
		},
		{
			name: "negative_1",
			mck: func(m *mocks.AdminProvider) {
				m.On("DeleteAdmin", mock.Anything, "awdgresbh").Return(false, storage.ErrUserNotFound)
			},
			args: args{
				login: "awdgresbh",
				key:   keyAdmin,
			},
			wantRes: false,
			wantErr: constants.ErrInvalidCredentials,
		},
		{
			name: "negative_2",
			mck: func(m *mocks.AdminProvider) {
				m.On("DeleteAdmin", mock.Anything, "awdgresbh").Return(false, errors.ErrUnsupported)
			},
			args: args{
				login: "awdgresbh",
				key:   keyAdmin,
			},
			wantRes: false,
			wantErr: constants.ErrInternalErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := mocks.NewAdminProvider(t)
			tt.mck(sqlite)

			s := &Auth{
				log:         slog.With(slog.String("service", "auth")),
				admProvider: sqlite,
			}

			gotRes, err := s.DeleteAdmin(context.Background(), tt.args.login, tt.args.key)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("DeleteAdmin() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestAuth_LoginUser(t *testing.T) {
	type mck func(m *mocks.UserProvider)

	type args struct {
		ctx      context.Context
		login    string
		password string
		appID    int32
	}
	tests := []struct {
		name      string
		mck       mck
		args      args
		wantToken string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Auth{
				log:         tt.fields.log,
				usrProvider: tt.fields.usrProvider,
				usrSaver:    tt.fields.usrSaver,
				appProvider: tt.fields.appProvider,
				admProvider: tt.fields.admProvider,
				tokenTTL:    tt.fields.tokenTTL,
			}
			gotToken, err := s.LoginUser(tt.args.ctx, tt.args.login, tt.args.password, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("LoginUser() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestAuth_RegisterNewUser(t *testing.T) {

	type mck func(m *mocks.UserProvider)
	type args struct {
		login    string
		password string
		appid    int32
	}
	tests := []struct {
		name       string
		mck        mck
		args       args
		wantUserid int64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Auth{
				log:         tt.fields.log,
				usrProvider: tt.fields.usrProvider,
				usrSaver:    tt.fields.usrSaver,
				appProvider: tt.fields.appProvider,
				admProvider: tt.fields.admProvider,
				tokenTTL:    tt.fields.tokenTTL,
			}
			gotUserid, err := s.RegisterNewUser(tt.args.ctx, tt.args.login, tt.args.password, tt.args.appid)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserid != tt.wantUserid {
				t.Errorf("RegisterNewUser() gotUserid = %v, want %v", gotUserid, tt.wantUserid)
			}
		})
	}
}
