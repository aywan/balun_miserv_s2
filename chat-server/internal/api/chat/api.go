package chat

import (
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"go.uber.org/zap"
)

type Implementation struct {
	desc.UnimplementedChatV1Server

	log      *zap.Logger
	chatServ service.Chat
}

func New(log *zap.Logger, chatServ service.Chat) *Implementation {
	return &Implementation{
		log:      log,
		chatServ: chatServ,
	}
}
