package grpc

import (
	"fmt"
	serverAPI "github.com/MorZLE/auth/internal/grpc/auth"
	"google.golang.org/grpc"

	"log/slog"
	"net"
)

func NewApp(log *slog.Logger, port int, authservice serverAPI.Auth, authAdmin serverAPI.AuthAdmin) *App {
	grpcServer := grpc.NewServer()

	serverAPI.RegisterServerAPI(grpcServer, authservice, authAdmin)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: grpcServer,
	}
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc.app.Run"
	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("running grpc server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpc.app.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
