package model

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int64         `db:"id"`
	CreatedAt time.Time     `db:"created_at"`
	UserID    sql.NullInt64 `db:"user_id"`
	MsgType   MessageType   `db:"type"`
	Text      string        `db:"text"`
}

type MessageList []Message
