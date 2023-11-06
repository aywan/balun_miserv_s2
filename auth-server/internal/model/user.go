package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Data      UserData     `db:""`
}

type UserData struct {
	Name         string `db:"name"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Role         int32  `db:"role"`
}
