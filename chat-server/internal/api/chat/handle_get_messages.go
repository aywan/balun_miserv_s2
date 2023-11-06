package chat

import (
	"context"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/converter"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"go.uber.org/zap"
)

func (s *Implementation) GetMessages(ctx context.Context, req *desc.GetMessagesRequest) (*desc.MessageListResponse, error) {
	result, err := s.chatServ.GetMessages(ctx, dto.GetMessagesDTO{
		ChatID:          req.ChatId,
		Limit:           req.Limit,
		AfterMessageId:  req.AfterMessageId,
		BeforeMessageId: req.BeforeMessageId,
	})
	if err != nil {
		s.log.Error("failed to get message list", zap.Error(err))
		return nil, err
	}

	items := make([]*desc.Message, 0, len(result.Items))
	for _, m := range result.Items {
		items = append(items, converter.ModelToGRPCMessage(m))
	}

	return &desc.MessageListResponse{
		Items:   items,
		HasNext: result.HasNext,
		NextId:  result.NextId,
	}, nil
}
