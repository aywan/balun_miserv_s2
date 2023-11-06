package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	userApi "github.com/aywan/balun_miserv_s2/auth-server/internal/api/user"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/app"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/config"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository"
	auditRepo "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/audit"
	userRepo "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/security"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/service"
	userService "github.com/aywan/balun_miserv_s2/auth-server/internal/service/user"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"github.com/aywan/balun_miserv_s2/shared/lib/logger"
	"github.com/aywan/balun_miserv_s2/shared/lib/runutil"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"github.com/fatih/color"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.MustNew(cfg.Mode)
	defer runutil.IgnoreErr(log.Sync)

	fxApp := fx.New(
		fx.Provide(
			func() *zap.Logger { return log },
			func() *db.Config { return &cfg.Db },
			func() *config.Server { return &cfg.Server },
			fx.Annotate(app.NewDb, fx.As(new(db.TxManager)), fx.As(new(db.DB))),
			fx.Annotate(userRepo.New, fx.As(new(repository.User))),
			fx.Annotate(auditRepo.New, fx.As(new(repository.Audit))),
			fx.Annotate(userService.New, fx.As(new(service.User))),
			fx.Annotate(app.NewGrpc, fx.As(new(grpc.ServiceRegistrar))),
			fx.Annotate(userApi.New, fx.As(new(desc.UserV1Server))),
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(
			func(server grpc.ServiceRegistrar) {},
			desc.RegisterUserV1Server,
		),
	)

	go func() {
		time.Sleep(time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := testConnectToServer(ctx, cfg.Server.Listen, log)
		if err != nil {
			log.Error("failed to test", zap.Error(err))
		}
		time.Sleep(time.Second)
		_ = fxApp.Stop(context.Background())
	}()

	fxApp.Run()
}

func testConnectToServer(ctx context.Context, serverAddr string, log *zap.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer runutil.LogOnError(conn.Close, log, "close client connection")

	c := desc.NewUserV1Client(conn)

	password := security.CreatePassword(6)
	rndEmail := fmt.Sprintf("some+%d@somemail.ru", rand.Int()) // #nosec G404 -- allow.
	createRsp, err := c.Create(ctx, &desc.CreateRequest{
		User: &desc.UserData{
			Name:  "SomeOne",
			Email: rndEmail,
			Role:  desc.UserRole_USER,
		},
		Credentials: &desc.UserCredentials{
			Password:        password,
			PasswordConfirm: password,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	getRsp, err := c.Get(ctx, &desc.UserIdRequest{Id: createRsp.Id})
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	log.Info(color.GreenString("id: %d", getRsp.GetId()))
	log.Info(color.GreenString("data: %+v", getRsp))

	rndEmail = fmt.Sprintf("new-some+%d@somemail.ru", rand.Int()) // #nosec G404 -- allow.
	_, err = c.Update(ctx, &desc.UpdateRequest{
		Id:    createRsp.Id,
		Email: wrapperspb.String(rndEmail),
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	_, err = c.Delete(ctx, &desc.UserIdRequest{Id: createRsp.Id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	_, err = c.Get(ctx, &desc.UserIdRequest{Id: createRsp.Id})
	if err == nil {
		return fmt.Errorf("user found after delete")
	}

	return nil
}
