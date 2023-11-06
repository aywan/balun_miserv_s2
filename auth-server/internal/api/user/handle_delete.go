package user

import (
	"context"

	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Delete(ctx context.Context, req *desc.UserIdRequest) (*emptypb.Empty, error) {
	err := s.userServ.Delete(ctx, req.Id)
	if err != nil {
		s.log.Error("failed to delete user", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
