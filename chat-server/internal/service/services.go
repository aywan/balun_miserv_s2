package service

import (
	"context"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
)

type Chat interface {
	Create(ctx context.Context, dto dto.NewChatDTO) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, data dto.SendMessageDTO) error
	GetMessages(ctx context.Context, req dto.GetMessagesDTO) (dto.MessagesResultDTO, error)
}
