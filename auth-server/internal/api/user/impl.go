package user

import (
	"github.com/aywan/balun_miserv_s2/auth-server/internal/service"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"go.uber.org/zap"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	log      *zap.Logger
	userServ service.User
}

func New(log *zap.Logger, userServ service.User) *Implementation {
	return &Implementation{
		log:      log,
		userServ: userServ,
	}
}
