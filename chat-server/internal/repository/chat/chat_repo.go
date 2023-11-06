package chat

import (
	"context"
	"fmt"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/repository"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/repository/chat/dto"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"go.uber.org/zap"
)

type Repo struct {
	db  db.DB
	log *zap.Logger
}

var _ repository.Chat = (*Repo)(nil)

func New(db db.DB, log *zap.Logger) *Repo {
	return &Repo{
		db:  db,
		log: log,
	}
}

func (r *Repo) CreateChat(ctx context.Context, data dto.CreateChatDTO) (int64, error) {
	builder := sq.Insert(tableChat).
		PlaceholderFormat(sq.Dollar).
		Columns(colChatOwnerId, colChatName).
		Values(data.OwnerID, data.Name).
		Suffix("RETURNING " + colChatId)

	query, err := db.BuildQuery("chat.create_chat", builder)
	if err != nil {
		return 0, err
	}

	var chatID int64
	err = r.db.QueryRowContext(ctx, query).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (r *Repo) CreateMessage(ctx context.Context, data dto.CreateMessageDTO) (int64, error) {
	builder := sq.Insert(tableMessage).
		PlaceholderFormat(sq.Dollar).
		Columns(colMessageChatId, colMessageUserId, colMessageType, colMessageText).
		Values(data.ChatId, data.UserId, data.Type, data.Text).
		Suffix("RETURNING " + colMessageId)

	query, err := db.BuildQuery("chat.new_message", builder)
	if err != nil {
		return 0, err
	}

	var msgID int64
	err = r.db.QueryRowContext(ctx, query).Scan(&msgID)
	if err != nil {
		return 0, err
	}

	return msgID, nil
}

func (r *Repo) UpdateChatLastMessage(ctx context.Context, chatID int64, lastMessageID int64) error {
	builder := sq.Insert(tableChatMessage+" as t").
		PlaceholderFormat(sq.Dollar).
		Columns(colChatMessageChatId, colChatMessageLastMessageId).
		Values(chatID, lastMessageID).
		Suffix(fmt.Sprintf(
			"ON CONFLICT (%[1]s) DO UPDATE SET %[2]s=greatest(excluded.%[2]s, t.%[2]s)",
			colChatMessageChatId,
			colChatMessageLastMessageId,
		))

	query, err := db.BuildQuery("chat.update_chat_last_message", builder)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query)
	return err
}

func (r *Repo) AddUsersToChat(ctx context.Context, chatID int64, users []int64, lastMessageID int64) error {
	builder := sq.Insert(tableChatUser).
		PlaceholderFormat(sq.Dollar).
		Columns(colChatUserChatId, colChatUserUserId, colChatUserLastMessageId)
	for _, userID := range users {
		builder = builder.Values(chatID, userID, lastMessageID)
	}

	query, err := db.BuildQuery("chat.add_users_to_chat", builder)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query)

	return err
}

func (r *Repo) DeleteChat(ctx context.Context, chatID int64) error {
	builder := sq.
		Update(tableChat).
		PlaceholderFormat(sq.Dollar).
		Set(colChatDeletedAt, sq.Expr("now()")).
		Set(colChatUpdatedAt, sq.Expr("now()")).
		Where(sq.Eq{colChatId: chatID})

	query, err := db.BuildQuery("chat.deleteChat", builder)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query)

	return err
}

func (r *Repo) GetMessageBefore(ctx context.Context, chatId int64, messageId int64, limit uint64) (model.MessageList, error) {
	builder := sq.
		Select(colMessageId, colMessageCreatedAt, colMessageUserId, colMessageType, colMessageText).
		From(tableMessage).
		PlaceholderFormat(sq.Dollar).
		Limit(limit).
		Where(sq.Lt{colMessageId: messageId}, sq.Eq{colMessageChatId: chatId}).
		OrderBy("id DESC")

	query, err := db.BuildQuery("chat.get_before_messages", builder)
	if err != nil {
		return nil, err
	}

	var messageList model.MessageList

	err = r.db.ScanAllContext(ctx, &messageList, query)
	if err != nil {
		return nil, err
	}

	slices.Reverse(messageList)

	return messageList, nil
}

func (r *Repo) GetMessageAfter(ctx context.Context, chatId int64, messageId int64, limit uint64) (model.MessageList, error) {
	builder := sq.
		Select(colMessageId, colMessageCreatedAt, colMessageUserId, colMessageType, colMessageText).
		From(tableMessage).
		PlaceholderFormat(sq.Dollar).
		Limit(limit).
		Where(sq.Gt{colMessageId: messageId}, sq.Eq{colMessageChatId: chatId}).
		OrderBy("id ASC")

	query, err := db.BuildQuery("chat.get_after_messages", builder)
	if err != nil {
		return nil, err
	}

	var messageList model.MessageList

	err = r.db.ScanAllContext(ctx, &messageList, query)
	if err != nil {
		return nil, err
	}

	return messageList, nil
}
