package user

import (
	"context"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/converter"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"go.uber.org/zap"
)

func (s *Implementation) Get(ctx context.Context, req *desc.UserIdRequest) (*desc.UserResponse, error) {
	userInst, err := s.userServ.Get(ctx, req.GetId())
	if err != nil {
		s.log.With(zap.Error(err)).Error("error getting user")
		return nil, err
	}

	rsp := converter.UserToGrpcUserResponse(userInst)

	return rsp, nil
}
