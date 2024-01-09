package app

import (
	grpcserver "github.com/MorZLE/auth/internal/app/grpc"
	"github.com/MorZLE/auth/internal/config"
	"github.com/MorZLE/auth/internal/controller/rest"
	"github.com/MorZLE/auth/internal/service"
	"github.com/MorZLE/auth/internal/storage/sqlite"
	"log/slog"
)

func NewApp(log *slog.Logger, cfg *config.Config) *App {

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		panic(err)
	}
	authservice := service.NewAuth(log, storage, storage, storage, storage, cfg.GRPC.Timeout)

	grpcApp := grpcserver.NewGRPC(log, cfg.GRPC.Port, authservice, authservice)

	restAPI := rest.NewHandler(log, authservice, authservice, cfg.Rest.Port, cfg.Rest.Timeout)

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
