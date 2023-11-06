package dto

import (
	"database/sql"
)

type UpdateDTO struct {
	Name         sql.NullString
	Email        sql.NullString
	PasswordHash sql.NullString
	Role         sql.NullInt32
}
