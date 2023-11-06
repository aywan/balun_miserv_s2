package converter

import (
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
)

func GRPCCreateChatReqToServiceDTO(req *desc.CreateRequest) dto.NewChatDTO {
	return dto.NewChatDTO{
		OwnerID: req.OwnerId,
		Name:    req.Name,
		Users:   req.Users,
	}
}
