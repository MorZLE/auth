package app

import (
	grpcserver "github.com/MorZLE/auth/internal/app/grpc"
	"github.com/MorZLE/auth/internal/controller/rest"
	"github.com/MorZLE/auth/internal/service"
	"github.com/MorZLE/auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

func NewApp(log *slog.Logger, port int, strPath string, ttl time.Duration) *App {

	storage, err := sqlite.NewStorage(strPath)
	if err != nil {
		panic(err)
	}
	authservice := service.NewAuth(log, storage, storage, storage, storage, ttl)
	grpcApp := grpcserver.NewApp(log, port, authservice, authservice)

	restAPI := rest.NewHandler(log, authservice, authservice)

	if err != nil {
		panic(err)
	}
	return &App{
		GRPCSrv: grpcApp,
		RESTapi: restAPI,
	}
}

type App struct {
	GRPCSrv *grpcserver.App
	RESTapi *rest.Handler
}
