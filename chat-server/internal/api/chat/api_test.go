package chat

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/mocks"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestApi_Create(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	req := &desc.CreateRequest{
		Users:   []int64{1, 2, 3},
		OwnerId: 1,
		Name:    gofakeit.BeerStyle(),
	}
	chatId := gofakeit.Int64()

	chatService := mocks.NewMockChat(t)
	chatService.EXPECT().
		Create(ctx, dto.NewChatDTO{
			OwnerID: req.OwnerId,
			Name:    req.Name,
			Users:   req.Users,
		}).
		Return(chatId, nil)

	api := New(log, chatService)

	rsp, err := api.Create(ctx, req)
	require.NoError(t, err)
	require.Equal(t, chatId, rsp.Id)
}

func TestApi_Delete(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	req := &desc.ChatIdRequest{
		Id: gofakeit.Int64(),
	}

	chatService := mocks.NewMockChat(t)
	chatService.EXPECT().
		Delete(ctx, req.Id).
		Return(nil)

	api := New(log, chatService)

	_, err = api.Delete(ctx, req)
	require.NoError(t, err)
}

func TestApi_GetMessages(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	req := &desc.GetMessagesRequest{
		ChatId:          gofakeit.Int64(),
		Limit:           gofakeit.Int64(),
		AfterMessageId:  gofakeit.Int64(),
		BeforeMessageId: gofakeit.Int64(),
	}

	messages := dto.MessagesResultDTO{
		Items: model.MessageList{
			{
				ID:        gofakeit.Int64(),
				CreatedAt: gofakeit.Date(),
				UserID:    sql.NullInt64{},
				MsgType:   model.MsgTypeSystem,
				Text:      gofakeit.BeerStyle(),
			},
			{
				ID:        gofakeit.Int64(),
				CreatedAt: gofakeit.Date(),
				UserID:    sql.NullInt64{gofakeit.Int64(), true},
				MsgType:   model.MsgTypeUser,
				Text:      gofakeit.BeerAlcohol() + " " + gofakeit.BeerIbu(),
			},
		},
		HasNext: gofakeit.Bool(),
		NextId:  gofakeit.Int64(),
	}

	chatService := mocks.NewMockChat(t)
	chatService.EXPECT().
		GetMessages(ctx, dto.GetMessagesDTO{
			ChatID:          req.ChatId,
			Limit:           req.Limit,
			AfterMessageId:  req.AfterMessageId,
			BeforeMessageId: req.BeforeMessageId,
		}).
		Return(messages, nil)

	api := New(log, chatService)

	rsp, err := api.GetMessages(ctx, req)
	require.NoError(t, err)

	require.Equal(t, messages.HasNext, rsp.HasNext)
	require.Equal(t, messages.NextId, rsp.NextId)

	require.Equal(t, messages.Items[0].ID, rsp.Items[0].Id)
	require.Equal(t, messages.Items[0].CreatedAt, rsp.Items[0].CreatedAt.AsTime())
	require.Nil(t, rsp.Items[0].UserId)
	require.Equal(t, desc.MessageType_SYSTEM, rsp.Items[0].Type)
	require.Equal(t, messages.Items[0].Text, rsp.Items[0].Text)

	require.Equal(t, messages.Items[1].ID, rsp.Items[1].Id)
	require.Equal(t, messages.Items[1].CreatedAt, rsp.Items[1].CreatedAt.AsTime())
	require.Equal(t, messages.Items[1].UserID.Int64, *rsp.Items[1].UserId)
	require.Equal(t, desc.MessageType_USER, rsp.Items[1].Type)
	require.Equal(t, messages.Items[1].Text, rsp.Items[1].Text)
}

func TestApi_SendMessage(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	req := &desc.SendMessageRequest{
		ChatId: gofakeit.Int64(),
		Type:   desc.MessageType_USER,
		UserId: gofakeit.Int64(),
		Text:   gofakeit.BeerStyle(),
	}

	chatService := mocks.NewMockChat(t)
	chatService.EXPECT().
		SendMessage(ctx, dto.SendMessageDTO{
			ChatID: req.ChatId,
			UserID: req.UserId,
			Text:   req.Text,
		}).
		Return(nil)

	api := New(log, chatService)

	_, err = api.SendMessage(ctx, req)
	require.NoError(t, err)
}
