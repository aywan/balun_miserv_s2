package chat

import (
	"context"

	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Delete(ctx context.Context, req *desc.ChatIdRequest) (*emptypb.Empty, error) {
	err := s.chatServ.Delete(ctx, req.Id)
	if err != nil {
		s.log.Error("failed to delete chat", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
