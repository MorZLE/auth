package app

import (
	grpcserver "github.com/MorZLE/auth/internal/app/grpc"
	"log/slog"
	"time"
)

func NewApp(log *slog.Logger, port int, strPath string, ttl time.Duration) *App {
	grpcApp := grpcserver.NewApp(log, port)

	return &App{
		GRPCSrv: grpcApp,
	}
}

type App struct {
	GRPCSrv *grpcserver.App
}
