package chat

import (
	"context"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/converter"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"go.uber.org/zap"
)

func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	data := converter.GRPCCreateChatReqToServiceDTO(req)
	chatID, err := s.chatServ.Create(ctx, data)
	if err != nil {
		s.log.Error("failed to create chat", zap.Error(err))
		return nil, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}
