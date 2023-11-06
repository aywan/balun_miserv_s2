package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/api/chat"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/app"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/config"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/repository"
	chatRepo "github.com/aywan/balun_miserv_s2/chat-server/internal/repository/chat"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service"
	chatService "github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"github.com/aywan/balun_miserv_s2/shared/lib/logger"
	"github.com/aywan/balun_miserv_s2/shared/lib/runutil"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
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
			fx.Annotate(chatRepo.New, fx.As(new(repository.Chat))),
			fx.Annotate(chatService.New, fx.As(new(service.Chat))),
			fx.Annotate(app.NewGrpc, fx.As(new(grpc.ServiceRegistrar))),
			fx.Annotate(chat.New, fx.As(new(desc.ChatV1Server))),
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(
			func(server grpc.ServiceRegistrar) {},
			desc.RegisterChatV1Server,
		),
	)

	go func() {
		time.Sleep(time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := testServer(ctx, cfg.Server.Listen, log)
		if err != nil {
			log.Error("failed to test", zap.Error(err))
		}
	}()

	fxApp.Run()
}

func testServer(ctx context.Context, addr string, log *zap.Logger) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer runutil.LogOnError(conn.Close, log, "error close client connection")

	c := desc.NewChatV1Client(conn)

	createReq, err := c.Create(ctx, &desc.CreateRequest{
		Users:   []int64{1, 2, 3},
		OwnerId: 1,
		Name:    "The chat",
	})
	if err != nil {
		return fmt.Errorf("failed to create chat, %w", err)
	}

	_, err = c.SendMessage(ctx, &desc.SendMessageRequest{
		ChatId: createReq.Id,
		UserId: 2,
		Text:   "text from two",
		Type:   desc.MessageType_USER,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	_, err = c.SendMessage(ctx, &desc.SendMessageRequest{
		ChatId: createReq.Id,
		UserId: 1,
		Text:   "text from one",
		Type:   desc.MessageType_USER,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	msgs, err := c.GetMessages(ctx, &desc.GetMessagesRequest{
		ChatId:          createReq.Id,
		Limit:           20,
		AfterMessageId:  0,
		BeforeMessageId: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}
	log.Info("get messages", zap.Any("messages", msgs.Items))

	_, err = c.Delete(ctx, &desc.ChatIdRequest{Id: createReq.Id})
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}
