package audit

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/audit/dto"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
)

const table = "\"audit\""

const (
	columnId          = "id"
	columnCreatedAt   = "created_at"
	columnCreatorID   = "creator_id"
	columnReference   = "reference"
	columnReferenceID = "reference_id"
	columnAction      = "action"
)

type Audit struct {
	db db.DB
}

func New(db db.DB) *Audit {
	return &Audit{db: db}
}

func (a *Audit) Insert(ctx context.Context, data dto.InsertDTO) (int64, error) {

	builder := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns(columnCreatorID, columnReference, columnReferenceID, columnAction).
		Values(data.CreatorId, data.Reference, data.ReferenceID, data.Action).
		Suffix("RETURNING " + columnId)

	query, err := db.BuildQuery("audit.insert", builder)
	if err != nil {
		return 0, err
	}

	var id int64
	err = a.db.QueryRowContext(ctx, query).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil

}
