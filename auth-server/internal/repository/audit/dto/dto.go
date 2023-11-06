package dto

import (
	"database/sql"
)

type InsertDTO struct {
	CreatorId   sql.NullInt64
	Reference   string
	ReferenceID int64
	Action      string
}
