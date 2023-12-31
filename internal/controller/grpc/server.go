package grpc

import (
	"context"
	"errors"
	"github.com/MorZLE/auth/internal/controller"
	"github.com/MorZLE/auth/internal/domain/cerror"
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth      controller.Auth
	authAdmin controller.AuthAdmin
}

func RegisterServerAPI(gRPC *grpc.Server, auth controller.Auth, authAdmin controller.AuthAdmin) {
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
		if errors.Is(err, cerror.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "login not found")
		}
		return nil, status.Error(codes.Internal, "internal cerror")
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
		if errors.Is(err, cerror.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal cerror")
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {

	userID := req.GetUserId()
	appID := req.GetAppId()
	if userID == emptyValue || appID == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "userID empty")
	}

	res, err := s.auth.CheckIsAdmin(ctx, userID, appID)
	if err != nil {
		if errors.Is(err, cerror.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal cerror")
	}

	return &authv1.IsAdminResponse{
		IsAdmin: true,
		Lvl:     res.Lvl,
	}, nil
}

func (s *serverAPI) CreateAdmin(ctx context.Context, req *authv1.CreateAdminRequest) (*authv1.CreateAdminResponse, error) {
	login := req.GetLogin()
	lvl := req.GetLvl()
	key := req.GetKey()
	appID := req.GetAppId()

	if login == "" || lvl == emptyValue || appID == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	userid, err := s.authAdmin.CreateAdmin(ctx, login, lvl, key, appID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotRights) {
			return nil, status.Error(codes.Aborted, "not enough rights")
		}
		if errors.Is(err, cerror.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal cerror")
	}
	return &authv1.CreateAdminResponse{UserId: userid}, nil
}

func (s *serverAPI) DeleteAdmin(ctx context.Context, req *authv1.DeleteAdminRequest) (*authv1.DeleteAdminResponse, error) {
	login := req.GetLogin()
	key := req.GetKey()

	if login == "" || key == "" {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	res, err := s.authAdmin.DeleteAdmin(ctx, login, key)
	if err != nil {
		if errors.Is(err, cerror.ErrNotRights) {
			return nil, status.Error(codes.NotFound, "user not admin")
		}
		return nil, status.Error(codes.Internal, "internal cerror")
	}
	return &authv1.DeleteAdminResponse{Result: res}, nil
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
		if errors.Is(err, cerror.ErrNotRights) {
			return nil, status.Error(codes.NotFound, "user not admin")
		}
		return nil, status.Error(codes.Internal, "internal cerror")
	}
	return &authv1.AddAppResponse{AppId: appID}, nil

}
