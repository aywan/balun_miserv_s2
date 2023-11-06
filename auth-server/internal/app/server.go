package app

import (
	"context"
	"net"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/config"
	"github.com/aywan/balun_miserv_s2/shared/lib/runutil"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGrpc(lc fx.Lifecycle, cfg *config.Server, log *zap.Logger) *grpc.Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			serverPort, err := net.Listen("tcp", cfg.Listen)
			if err != nil {
				return nil
			}

			go runutil.LogOnError(
				func() error {
					return grpcServer.Serve(serverPort)
				},
				log,
				"grpc server start error",
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	return grpcServer
}
