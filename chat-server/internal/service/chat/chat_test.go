package chat

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
	repoDto "github.com/aywan/balun_miserv_s2/chat-server/internal/repository/chat/dto"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/repository/mocks"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/service/chat/dto"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestService_Create(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	newChatDTO := dto.NewChatDTO{
		OwnerID: gofakeit.Int64(),
		Name:    gofakeit.BeerStyle(),
		Users:   []int64{gofakeit.Int64(), gofakeit.Int64()},
	}
	chatId := gofakeit.Int64()
	msgId := gofakeit.Int64()

	chatRepoMock := mocks.NewMockChat(t)
	chatRepoMock.EXPECT().
		CreateChat(ctx, repoDto.CreateChatDTO{
			OwnerID: newChatDTO.OwnerID,
			Name:    newChatDTO.Name,
		}).
		Return(chatId, nil)

	chatRepoMock.EXPECT().
		CreateMessage(ctx, repoDto.CreateMessageDTO{
			ChatId: chatId,
			UserId: sql.NullInt64{},
			Text:   "start new chat",
			Type:   model.MsgTypeSystem,
		}).
		Return(msgId, nil)

	chatRepoMock.EXPECT().
		UpdateChatLastMessage(ctx, chatId, msgId).
		Return(nil)

	chatRepoMock.EXPECT().
		AddUsersToChat(ctx, chatId, newChatDTO.Users, msgId).
		Return(nil)

	txManager := db.NewTestTxManager(t)

	service := New(log, txManager, chatRepoMock)

	actualChatId, err := service.Create(ctx, newChatDTO)
	require.NoError(t, err)
	require.Equal(t, chatId, actualChatId)
}

func TestService_Delete(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	chatId := gofakeit.Int64()

	chatRepoMock := mocks.NewMockChat(t)
	chatRepoMock.EXPECT().
		DeleteChat(ctx, chatId).
		Return(nil)

	txManager := db.NewTestTxManager(t)

	service := New(log, txManager, chatRepoMock)

	err = service.Delete(ctx, chatId)
	require.NoError(t, err)
}

func TestService_SendMessage(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	msgDto := dto.SendMessageDTO{
		ChatID: gofakeit.Int64(),
		UserID: gofakeit.Int64(),
		Text:   gofakeit.BeerName(),
	}

	msgId := gofakeit.Int64()

	chatRepoMock := mocks.NewMockChat(t)
	chatRepoMock.EXPECT().
		CreateMessage(ctx, repoDto.CreateMessageDTO{
			ChatId: msgDto.ChatID,
			UserId: sql.NullInt64{msgDto.UserID, true},
			Text:   msgDto.Text,
			Type:   model.MsgTypeUser,
		}).
		Return(msgId, nil)

	chatRepoMock.EXPECT().
		UpdateChatLastMessage(ctx, msgDto.ChatID, msgId).
		Return(nil)

	txManager := db.NewTestTxManager(t)

	service := New(log, txManager, chatRepoMock)

	err = service.SendMessage(ctx, msgDto)
	require.NoError(t, err)
}

func TestService_GetMessages(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()

	messages := model.MessageList{
		{
			ID:        1,
			CreatedAt: time.Time{},
			UserID:    sql.NullInt64{},
			MsgType:   model.MsgTypeSystem,
			Text:      "System",
		},
		{
			ID:        2,
			CreatedAt: time.Time{},
			UserID:    sql.NullInt64{22, true},
			MsgType:   model.MsgTypeUser,
			Text:      "User",
		},
	}

	chatRepoMock := mocks.NewMockChat(t)
	chatRepoMock.EXPECT().
		GetMessageBefore(ctx, mock.Anything, mock.Anything, mock.Anything).
		Return(messages, nil).
		Maybe()
	chatRepoMock.EXPECT().
		GetMessageAfter(ctx, mock.Anything, mock.Anything, mock.Anything).
		Return(messages, nil).
		Maybe()
	txManager := db.NewTestTxManager(t)
	service := New(log, txManager, chatRepoMock)

	cases := []struct {
		Name          string
		Req           dto.GetMessagesDTO
		ExpectIds     []int64
		ExpectHasNext bool
		ExpectNextId  int64
		ExpectedErr   error
	}{
		{
			Name: "after limit 2",
			Req: dto.GetMessagesDTO{
				Limit:           2,
				AfterMessageId:  1,
				BeforeMessageId: 0,
			},
			ExpectIds:     []int64{1, 2},
			ExpectHasNext: false,
			ExpectNextId:  0,
		},
		{
			Name: "after limit 1",
			Req: dto.GetMessagesDTO{
				Limit:           1,
				AfterMessageId:  1,
				BeforeMessageId: 0,
			},
			ExpectIds:     []int64{1},
			ExpectHasNext: true,
			ExpectNextId:  2,
		},
		{
			Name: "before limit 2",
			Req: dto.GetMessagesDTO{
				Limit:           2,
				AfterMessageId:  0,
				BeforeMessageId: 2,
			},
			ExpectIds:     []int64{1, 2},
			ExpectHasNext: false,
			ExpectNextId:  0,
		},
		{
			Name: "before limit 1",
			Req: dto.GetMessagesDTO{
				Limit:           1,
				AfterMessageId:  0,
				BeforeMessageId: 2,
			},
			ExpectIds:     []int64{2},
			ExpectHasNext: true,
			ExpectNextId:  1,
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.Name, func(t *testing.T) {
			t.Parallel()

			rsp, errActual := service.GetMessages(ctx, cc.Req)
			if cc.ExpectedErr != nil {
				require.Error(t, errActual)
				require.ErrorIs(t, errActual, cc.ExpectedErr)
				return
			}

			actualIds := make([]int64, 0, len(rsp.Items))
			for _, item := range rsp.Items {
				actualIds = append(actualIds, item.ID)
			}
			require.Equal(t, cc.ExpectIds, actualIds)
			require.Equal(t, cc.ExpectHasNext, rsp.HasNext)
			require.Equal(t, cc.ExpectNextId, rsp.NextId)
		})
	}
}
