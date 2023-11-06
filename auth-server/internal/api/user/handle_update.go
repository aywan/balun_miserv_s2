package user

import (
	"context"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/converter"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	updDto := converter.GrpcUpdateRequestToUpdateDTO(req)

	err := s.userServ.Update(ctx, req.Id, updDto)
	if err != nil {
		s.log.Error("failed to update user", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
