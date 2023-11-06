package repository

import (
	"context"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/repository/chat/dto"
)

type Chat interface {
	CreateChat(ctx context.Context, data dto.CreateChatDTO) (int64, error)
	CreateMessage(ctx context.Context, data dto.CreateMessageDTO) (int64, error)
	UpdateChatLastMessage(ctx context.Context, chatID int64, lastMessageID int64) error
	AddUsersToChat(ctx context.Context, chatID int64, users []int64, lastMessageID int64) error
	DeleteChat(ctx context.Context, chatID int64) error
	GetMessageBefore(ctx context.Context, chatId int64, messageId int64, limit uint64) (model.MessageList, error)
	GetMessageAfter(ctx context.Context, chatId int64, messageId int64, limit uint64) (model.MessageList, error)
}
