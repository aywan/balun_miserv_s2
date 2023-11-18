package dto

import (
	"database/sql"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
)

type InsertDTO struct {
	CreatorId   sql.NullInt64
	Reference   model.AuditReference
	ReferenceID int64
	Action      model.AuditAction
}
