package user

import (
	"context"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/converter"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"go.uber.org/zap"
)

func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userData := converter.GrpcCreateRequestToUserData(req)

	userId, err := s.userServ.Create(ctx, userData)
	if err != nil {
		s.log.Error("failed to insert user", zap.Error(err))
		return nil, err
	}

	return &desc.CreateResponse{Id: userId}, nil
}
