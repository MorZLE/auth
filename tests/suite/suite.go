package suite

import (
	"context"
	"github.com/MorZLE/auth/internal/config"
	authv1 "github.com/MorZLE/auth/internal/generate/grpc/gen/morzle.auth.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

const host = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authv1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {

	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("..\\config\\local.yaml")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	authClient, err := grpc.DialContext(context.Background(),
		net.JoinHostPort(host, strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authv1.NewAuthClient(authClient),
	}
}
