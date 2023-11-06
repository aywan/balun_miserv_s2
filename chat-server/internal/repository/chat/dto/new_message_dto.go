package dto

import (
	"database/sql"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
)

type CreateMessageDTO struct {
	ChatId int64
	UserId sql.NullInt64
	Text   string
	Type   model.MessageType
}
