package auth

import (
	"context"
	"errors"
	"github.com/MorZLE/auth/internal/constants"
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"github.com/MorZLE/auth/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	LoginUser(ctx context.Context, login string, password string, appID int32) (token string, err error)
	RegisterNewUser(ctx context.Context, login string, password string, appid int32) (userid int64, err error)
	CheckIsAdmin(ctx context.Context, userid int32) (bool, error)
}

type AuthAdmin interface {
	CreateAdmin(ctx context.Context, login string, lvl int32, key string) (userid int64, err error)
	DeleteAdmin(ctx context.Context, login string, key string) (userid int64, err error)
	AddApp(ctx context.Context, name, secret, key string) (userid int32, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth      Auth
	authAdmin AuthAdmin
}

func RegisterServerAPI(gRPC *grpc.Server, auth Auth, authAdmin AuthAdmin) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth, authAdmin: authAdmin})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	login := req.GetLogin()
	pswrd := req.GetPassword()
	numApp := req.GetAppId()

	if login == "" || numApp == emptyValue || pswrd == "" {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	token, err := s.auth.LoginUser(ctx, login, pswrd, numApp)
	if err != nil {
		if errors.Is(err, constants.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "internal error")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	login := req.GetLogin()
	pswrd := req.GetPassword()
	appid := req.GetAppId()

	if login == "" || pswrd == "" || appid == 0 {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	userID, err := s.auth.RegisterNewUser(ctx, login, pswrd, appid)
	if err != nil {
		if errors.Is(err, constants.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {

	userID := req.UserId
	if userID == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "userID exist")
	}

	flag, err := s.auth.CheckIsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "internal error")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.IsAdminResponse{
		IsAdmin: flag,
	}, nil
}

func (s *serverAPI) CreateAdmin(ctx context.Context, req *authv1.CreateAdminRequest) (*authv1.CreateAdminResponse, error) {
	login := req.GetLogin()
	lvl := req.GetLvl()
	key := req.GetKey()

	if login == "" || lvl == 0 || key == "" {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	userid, err := s.authAdmin.CreateAdmin(ctx, login, lvl, key)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.CreateAdminResponse{UserId: userid}, nil
}

func (s *serverAPI) DeleteAdmin(ctx context.Context, req *authv1.DeleteAdminRequest) (*authv1.DeleteAdminResponse, error) {
	login := req.GetLogin()
	key := req.GetKey()

	if login == "" || key == "" {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	userid, err := s.authAdmin.DeleteAdmin(ctx, login, key)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.DeleteAdminResponse{UserId: userid}, nil
}

func (s *serverAPI) AddApp(ctx context.Context, req *authv1.AddAppRequest) (*authv1.AddAppResponse, error) {
	name := req.GetName()
	secret := req.GetSecret()
	key := req.GetKey()

	if name == "" || secret == "" || key == "" {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}
	appID, err := s.authAdmin.AddApp(ctx, name, secret, key)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.AddAppResponse{AppId: appID}, nil

}
