package grpc

import (
	"context"
	"errors"
	"github.com/MorZLE/auth/internal/controller/grpc/mocks"
	"github.com/MorZLE/auth/internal/domain/cerror"
	"github.com/MorZLE/auth/internal/domain/models"
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"testing"
)

func Test_serverAPI_Login(t *testing.T) {
	type mck func(m *mocks.Auth)

	type args struct {
		ctx context.Context
		req *authv1.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		mck     mck
		want    *authv1.LoginResponse
		wantErr error
	}{
		{
			name: "positive_1",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "test",
					Password: "test",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("LoginUser", context.Background(), "test", "test", int32(1)).Return("token", nil)
			},
			want: &authv1.LoginResponse{
				Token: "token",
			},
			wantErr: nil,
		},
		{
			name: "positive_2",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "12edsc231df",
					Password: "test",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("LoginUser", context.Background(), "12edsc231df", "test", int32(1)).Return("q3egersthg435h4wh5tjhr67ksazrtjnh54y", nil)
			},
			want: &authv1.LoginResponse{
				Token: "q3egersthg435h4wh5tjhr67ksazrtjnh54y",
			},
			wantErr: nil,
		},
		{
			name: "empty_login",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "",
					Password: "test",
					AppId:    1,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_password",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "test",
					Password: "",
					AppId:    1,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_appid",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "stfrh",
					Password: "sth",
					AppId:    0,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_data",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "",
					Password: "",
					AppId:    0,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "login_not_found",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "teset",
					Password: "teset",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("LoginUser", context.Background(), "teset", "teset", int32(1)).Return("", cerror.ErrInvalidCredentials)
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "login not found"),
		},
		{
			name: "login_not_exist",
			args: args{
				ctx: context.Background(),
				req: &authv1.LoginRequest{
					Login:    "teset",
					Password: "teset",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("LoginUser", context.Background(), "teset", "teset", int32(1)).Return("", errors.ErrUnsupported)
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal cerror"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serAuth := mocks.NewAuth(t)
			tt.mck(serAuth)
			s := &serverAPI{
				auth: serAuth,
			}
			got, err := s.Login(tt.args.ctx, tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Login() cerror = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serverAPI_Register(t *testing.T) {

	type mck func(m *mocks.Auth)
	type args struct {
		req *authv1.RegisterRequest
	}
	tests := []struct {
		name    string
		mck     mck
		args    args
		want    *authv1.RegisterResponse
		wantErr error
	}{
		{
			name: "positive_1",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "test",
					Password: "test",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("RegisterNewUser", context.Background(), "test", "test", int32(1)).Return(int64(1), nil)
			},
			want: &authv1.RegisterResponse{
				UserId: 1,
			},
			wantErr: nil,
		},
		{
			name: "positive_2",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "afsderbhersnb",
					Password: "tsanbsrtnrsnest",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("RegisterNewUser", context.Background(), "afsderbhersnb", "tsanbsrtnrsnest", int32(1)).Return(int64(2), nil)
			},
			want: &authv1.RegisterResponse{
				UserId: 2,
			},
			wantErr: nil,
		},
		{
			name: "empty_login",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "",
					Password: "tsanbsrtnrsnest",
					AppId:    1,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_Password",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "dzhdtjhrsjrstjh",
					Password: "",
					AppId:    1,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_AppId",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "dzhdtjhrsjrstjh",
					Password: "dzhdtjhrsjrstjh",
					AppId:    0,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "loginExist",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "dzhdtjhrsjrstjh",
					Password: "dzhdtjhrsjrstjh",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("RegisterNewUser", context.Background(), "dzhdtjhrsjrstjh", "dzhdtjhrsjrstjh", int32(1)).Return(int64(0), cerror.ErrUserExists)
			},
			want:    nil,
			wantErr: status.Error(codes.AlreadyExists, "user already exists"),
		},
		{
			name: "internal_error",
			args: args{
				req: &authv1.RegisterRequest{
					Login:    "dzhdtjhrsjrstjh",
					Password: "dzhdtjhrsjrstjh",
					AppId:    1,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("RegisterNewUser", context.Background(), "dzhdtjhrsjrstjh", "dzhdtjhrsjrstjh", int32(1)).Return(int64(0), errors.New("internal cerror"))
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal cerror"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewAuth(t)
			tt.mck(service)
			s := &serverAPI{
				auth: service,
			}
			got, err := s.Register(context.Background(), tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Register() cerror = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serverAPI_IsAdmin(t *testing.T) {
	type mck func(m *mocks.Auth)
	type args struct {
		req *authv1.IsAdminRequest
	}
	tests := []struct {
		name    string
		mck     mck
		args    args
		want    *authv1.IsAdminResponse
		wantErr error
	}{
		{
			name: "positive_1",
			mck: func(m *mocks.Auth) {
				m.On("CheckIsAdmin", context.Background(), int32(1), int32(1)).Return(models.Admin{Lvl: 1}, nil)
			},
			args: args{
				req: &authv1.IsAdminRequest{
					UserId: 1,
					AppId:  1,
				},
			},
			want: &authv1.IsAdminResponse{
				IsAdmin: true,
				Lvl:     1,
			},
			wantErr: nil,
		},
		{
			name: "positive_2",
			mck: func(m *mocks.Auth) {
				m.On("CheckIsAdmin", context.Background(), int32(2), int32(2)).Return(models.Admin{Lvl: 2}, nil)
			},
			args: args{
				req: &authv1.IsAdminRequest{
					UserId: 2,
					AppId:  2,
				},
			},
			want: &authv1.IsAdminResponse{
				IsAdmin: true,
				Lvl:     2,
			},

			wantErr: nil,
		},
		{
			name: "empty userID",
			args: args{
				req: &authv1.IsAdminRequest{
					UserId: 0,
					AppId:  2,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "userID empty"),
		},
		{
			name: "empty appID",
			args: args{
				req: &authv1.IsAdminRequest{
					UserId: 6,
					AppId:  0,
				},
			},
			mck:     func(m *mocks.Auth) {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "userID empty"),
		},
		{
			name: "user not exist",
			args: args{
				req: &authv1.IsAdminRequest{
					UserId: 6,
					AppId:  3,
				},
			},
			mck: func(m *mocks.Auth) {
				m.On("CheckIsAdmin", context.Background(), int32(6), int32(3)).Return(models.Admin{}, cerror.ErrInvalidCredentials)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, "user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewAuth(t)
			tt.mck(service)

			s := &serverAPI{
				auth: service,
			}
			got, err := s.IsAdmin(context.Background(), tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("IsAdmin() cerror = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsAdmin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serverAPI_CreateAdmin(t *testing.T) {
	type mck func(m *mocks.AuthAdmin)
	type args struct {
		req *authv1.CreateAdminRequest
	}
	tests := []struct {
		name    string
		mck     mck
		args    args
		want    *authv1.CreateAdminResponse
		wantErr error
	}{
		{
			name: "positive_1",
			mck: func(m *mocks.AuthAdmin) {
				m.On("CreateAdmin", context.Background(), "admin", int32(1), "key", int32(1)).Return(int64(1), nil)
			},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "admin",
					Lvl:   1,
					Key:   "key",
					AppId: 1,
				},
			},
			want: &authv1.CreateAdminResponse{
				UserId: 1,
			},
		},
		{
			name: "positive_2",
			mck: func(m *mocks.AuthAdmin) {
				m.On("CreateAdmin", context.Background(), "ZGzdbdfgber", int32(23), "Rbvzdb", int32(23)).Return(int64(1234), nil)
			},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "ZGzdbdfgber",
					Lvl:   23,
					Key:   "Rbvzdb",
					AppId: 23,
				},
			},
			want: &authv1.CreateAdminResponse{
				UserId: 1234,
			},
		},
		{
			name: "empty_login",
			mck:  func(m *mocks.AuthAdmin) {},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "",
					Lvl:   23,
					Key:   "Rbvzdb",
					AppId: 23,
				},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_lvl",
			mck:  func(m *mocks.AuthAdmin) {},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "se",
					Lvl:   0,
					Key:   "Rbvzdb",
					AppId: 23,
				},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty_appId",
			mck:  func(m *mocks.AuthAdmin) {},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "sefsef",
					Lvl:   23,
					Key:   "Rbvzdb",
					AppId: 0,
				},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "invalid_Key",
			mck: func(m *mocks.AuthAdmin) {
				m.On("CreateAdmin", context.Background(), "sefsef", int32(23), "", int32(43)).Return(int64(0), cerror.ErrNotRights)
			},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "sefsef",
					Lvl:   23,
					Key:   "",
					AppId: 43,
				},
			},
			want:    nil,
			wantErr: status.Error(codes.Aborted, "not enough rights"),
		},
		{
			name: "internal_error",
			mck: func(m *mocks.AuthAdmin) {
				m.On("CreateAdmin", context.Background(), "sefsef", int32(23), "", int32(890)).Return(int64(0), errors.New("cerror"))
			},
			args: args{
				req: &authv1.CreateAdminRequest{
					Login: "sefsef",
					Lvl:   23,
					Key:   "",
					AppId: 890,
				},
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal cerror"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewAuthAdmin(t)
			tt.mck(service)
			s := &serverAPI{
				authAdmin: service,
			}
			got, err := s.CreateAdmin(context.Background(), tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateAdmin() cerror = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAdmin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serverAPI_DeleteAdmin(t *testing.T) {
	type mck func(m *mocks.AuthAdmin)
	type args struct {
		req *authv1.DeleteAdminRequest
	}
	tests := []struct {
		name    string
		mck     mck
		args    args
		want    *authv1.DeleteAdminResponse
		wantErr error
	}{
		{
			name: "positive_1",
			mck: func(m *mocks.AuthAdmin) {
				m.On("DeleteAdmin", context.Background(), "sefsef", "wqrqwre").Return(true, nil)
			},
			args: args{
				req: &authv1.DeleteAdminRequest{
					Login: "sefsef",
					Key:   "wqrqwre",
				},
			},
			want: &authv1.DeleteAdminResponse{
				Result: true,
			},
		},
		{
			name: "positive_2",
			mck: func(m *mocks.AuthAdmin) {
				m.On("DeleteAdmin", context.Background(), "sesfsgdb43g", "argearg").Return(true, nil)
			},
			args: args{
				req: &authv1.DeleteAdminRequest{
					Login: "sesfsgdb43g",
					Key:   "argearg",
				},
			},
			want: &authv1.DeleteAdminResponse{
				Result: true,
			},
		},
		{
			name: "empty login",
			mck:  func(m *mocks.AuthAdmin) {},
			args: args{
				req: &authv1.DeleteAdminRequest{
					Login: "",
					Key:   "argearg",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "invalid key",
			mck: func(m *mocks.AuthAdmin) {
				m.On("DeleteAdmin", context.Background(), "awdvzvwe", "argearg").Return(false, cerror.ErrNotRights)
			},
			args: args{
				req: &authv1.DeleteAdminRequest{
					Login: "awdvzvwe",
					Key:   "argearg",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, "user not admin"),
		},
		{
			name: "internal cerror",
			mck: func(m *mocks.AuthAdmin) {
				m.On("DeleteAdmin", context.Background(), "awdvzvwe", "argearg").Return(false, errors.ErrUnsupported)
			},
			args: args{
				req: &authv1.DeleteAdminRequest{
					Login: "awdvzvwe",
					Key:   "argearg",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal cerror"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewAuthAdmin(t)
			tt.mck(service)
			s := &serverAPI{
				authAdmin: service,
			}
			got, err := s.DeleteAdmin(context.Background(), tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateAdmin() cerror = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAdmin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serverAPI_AddApp(t *testing.T) {
	type mck func(m *mocks.AuthAdmin)

	type args struct {
		req *authv1.AddAppRequest
	}
	tests := []struct {
		name    string
		args    args
		mck     mck
		want    *authv1.AddAppResponse
		wantErr error
	}{
		{
			name: "positive_1",
			mck: func(m *mocks.AuthAdmin) {
				m.On("AddApp", context.Background(), "sefsef", "wqrqwre", "sefsfe").Return(int32(1), nil)
			},
			args: args{
				req: &authv1.AddAppRequest{
					Name:   "sefsef",
					Secret: "wqrqwre",
					Key:    "sefsfe",
				},
			},
			want: &authv1.AddAppResponse{
				AppId: 1,
			},
		},
		{
			name: "positive_2",
			mck: func(m *mocks.AuthAdmin) {
				m.On("AddApp", context.Background(), "saZGasrgsd", "sdebdzbf", "sefsfe").Return(int32(1), nil)
			},
			args: args{
				req: &authv1.AddAppRequest{
					Name:   "saZGasrgsd",
					Secret: "sdebdzbf",
					Key:    "sefsfe",
				},
			},
			want: &authv1.AddAppResponse{
				AppId: 1,
			},
		},
		{
			name: "empty name",
			mck:  func(m *mocks.AuthAdmin) {},
			args: args{
				req: &authv1.AddAppRequest{
					Name:   "",
					Secret: "sdebdzbf",
					Key:    "sefsfe",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "empty secret",
			mck:  func(m *mocks.AuthAdmin) {},
			args: args{
				req: &authv1.AddAppRequest{
					Name:   "sefsef",
					Secret: "",
					Key:    "sefsfe",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "data not exist"),
		},
		{
			name: "negative key",
			mck: func(m *mocks.AuthAdmin) {
				m.On("AddApp", context.Background(), "sefsef", "wqrqwre", "sefsfe").Return(int32(0), cerror.ErrNotRights)
			},
			args: args{
				req: &authv1.AddAppRequest{
					Name:   "sefsef",
					Secret: "wqrqwre",
					Key:    "sefsfe",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, "user not admin"),
		},
		{
			name: "internal cerror",
			mck: func(m *mocks.AuthAdmin) {
				m.On("AddApp", context.Background(), "sefsef", "wqrqwre", "sefsfe").Return(int32(0), errors.New("internal cerror"))
			},
			args: args{
				req: &authv1.AddAppRequest{
					Name:   "sefsef",
					Secret: "wqrqwre",
					Key:    "sefsfe",
				},
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal cerror"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewAuthAdmin(t)
			tt.mck(service)

			s := &serverAPI{
				authAdmin: service,
			}
			got, err := s.AddApp(context.Background(), tt.args.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddApp() cerror = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddApp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
