package chat

import (
	"context"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.chatServ.SendMessage(ctx, dto.SendMessageDTO{
		ChatID: req.ChatId,
		UserID: req.UserId,
		Text:   req.Text,
	})

	if err != nil {
		s.log.Error("failed to send new message to chat", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
