package auth

import (
	"context"
	"errors"
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"github.com/MorZLE/auth/internal/service"
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
	RegisterNewUser(ctx context.Context, login string, password string) (userid int64, err error)
	CheckIsAdmin(ctx context.Context, userid int32) (bool, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func RegisterServerAPI(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
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
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "internal error")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	login := req.GetLogin()
	pswrd := req.GetPassword()

	if login == "" || pswrd == "" {
		return nil, status.Error(codes.InvalidArgument, "data not exist")
	}

	userID, err := s.auth.RegisterNewUser(ctx, login, pswrd)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "internal error")
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
