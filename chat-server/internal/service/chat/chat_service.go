package chat

import (
	"context"
	"database/sql"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/repository"
	repoDto "github.com/aywan/balun_miserv_s2/chat-server/internal/repository/chat/dto"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"go.uber.org/zap"
)

type Service struct {
	log       *zap.Logger
	txManager db.TxManager
	chatRepo  repository.Chat
}

var _ service.Chat = (*Service)(nil)

func New(log *zap.Logger, txManager db.TxManager, chatRepo repository.Chat) *Service {
	return &Service{
		log:       log,
		txManager: txManager,
		chatRepo:  chatRepo,
	}
}

func (s *Service) Create(ctx context.Context, data dto.NewChatDTO) (int64, error) {
	var outChatId int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		chatId, err := s.chatRepo.CreateChat(ctx, repoDto.CreateChatDTO{
			OwnerID: data.OwnerID,
			Name:    data.Name,
		})
		if err != nil {
			return err
		}

		msgId, err := s.chatRepo.CreateMessage(ctx, repoDto.CreateMessageDTO{
			ChatId: chatId,
			UserId: sql.NullInt64{},
			Text:   "start new chat",
			Type:   model.MsgTypeSystem,
		})
		if err != nil {
			return err
		}

		err = s.chatRepo.UpdateChatLastMessage(ctx, chatId, msgId)
		if err != nil {
			return err
		}

		err = s.chatRepo.AddUsersToChat(ctx, chatId, data.Users, msgId)
		if err != nil {
			return err
		}

		outChatId = chatId
		return nil
	})
	if err != nil {
		return 0, err
	}

	return outChatId, nil
}

func (s *Service) Delete(ctx context.Context, chatId int64) error {
	return s.chatRepo.DeleteChat(ctx, chatId)
}

func (s *Service) SendMessage(ctx context.Context, data dto.SendMessageDTO) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		msgId, err := s.chatRepo.CreateMessage(ctx, repoDto.CreateMessageDTO{
			ChatId: data.ChatID,
			UserId: sql.NullInt64{Valid: true, Int64: data.UserID},
			Text:   data.Text,
			Type:   model.MsgTypeUser,
		})
		if err != nil {
			return err
		}

		err = s.chatRepo.UpdateChatLastMessage(ctx, data.ChatID, msgId)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) GetMessages(ctx context.Context, req dto.GetMessagesDTO) (dto.MessagesResultDTO, error) {
	limit := uint64(req.Limit) + 1
	hasNext := false
	nextMessageId := int64(0)

	var mesg model.MessageList
	var err error
	if req.BeforeMessageId > 0 {
		mesg, err = s.chatRepo.GetMessageBefore(ctx, req.ChatID, req.BeforeMessageId, limit)
		if err != nil {
			return dto.MessagesResultDTO{}, err
		}
		if len(mesg) == int(limit) {
			hasNext = true
			nextMessageId = mesg[0].ID
			mesg = mesg[1:]
		}
	} else {
		mesg, err = s.chatRepo.GetMessageAfter(ctx, req.ChatID, req.AfterMessageId, limit)
		if err != nil {
			return dto.MessagesResultDTO{}, err
		}
		if len(mesg) == int(limit) {
			hasNext = true
			nextMessageId = mesg[len(mesg)-1].ID
			mesg = mesg[0 : len(mesg)-1]
		}
	}

	return dto.MessagesResultDTO{
		Items:   mesg,
		HasNext: hasNext,
		NextId:  nextMessageId,
	}, nil
}
