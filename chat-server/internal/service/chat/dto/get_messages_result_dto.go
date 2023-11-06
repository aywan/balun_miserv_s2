package dto

import (
	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
)

type MessagesResultDTO struct {
	Items   model.MessageList
	HasNext bool
	NextId  int64
}
