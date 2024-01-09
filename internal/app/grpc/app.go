package grpc

import (
	"fmt"
	"github.com/MorZLE/auth/internal/controller"
	serverAPI "github.com/MorZLE/auth/internal/controller/grpc"
	"google.golang.org/grpc"

	"log/slog"
	"net"
)

func NewGRPC(log *slog.Logger, port int, authservice controller.Auth, authAdmin controller.AuthAdmin) *App {
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
	const op = "controller.app.Run"
	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("running controller server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "controller.app.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
