package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
)

const table = "\"user\""

const (
	columnId           = "id"
	columnCreatedAt    = "created_at"
	columnUpdatedAt    = "updated_at"
	columnDeletedAt    = "deleted_at"
	columnName         = "name"
	columnEmail        = "email"
	columnPasswordHash = "password_hash"
	columnRole         = "role"
)

type Repository struct {
	db db.DB
}

var _ repository.User = (*Repository)(nil)

func New(db db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetNotDeleted(ctx context.Context, userId int64) (model.User, error) {
	builder := sq.
		Select("*").
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			columnId:        userId,
			columnDeletedAt: nil,
		})

	query, err := db.BuildQuery("user_repo.get", builder)
	if err != nil {
		return model.User{}, err
	}

	var user model.User

	err = r.db.ScanOneContext(ctx, &user, query)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, data model.UserData) (int64, error) {
	builder := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns(columnName, columnEmail, columnRole, columnPasswordHash).
		Values(data.Name, data.Email, data.Role, data.PasswordHash).
		Suffix("RETURNING " + columnId)

	query, err := db.BuildQuery("user_repo.create", builder)
	if err != nil {
		return 0, err
	}

	var userId int64
	err = r.db.QueryRowContext(ctx, query).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *Repository) Update(ctx context.Context, userId int64, data dto.UpdateDTO) error {
	builder := sq.
		Update(table).
		PlaceholderFormat(sq.Dollar).
		Set(columnUpdatedAt, sq.Expr("now()")).
		Where(sq.Eq{columnId: userId})

	if data.Name.Valid {
		builder = builder.Set(columnName, data.Name)
	}
	if data.Email.Valid {
		builder = builder.Set(columnEmail, data.Email)
	}
	if data.PasswordHash.Valid {
		builder = builder.Set(columnPasswordHash, data.PasswordHash)
	}
	if data.Role.Valid {
		builder = builder.Set(columnRole, data.Role)
	}

	query, err := db.BuildQuery("user_repo.update", builder)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query)

	return err
}

func (r *Repository) Delete(ctx context.Context, userId int64) error {
	builder := sq.
		Update(table).
		PlaceholderFormat(sq.Dollar).
		Set(columnDeletedAt, sq.Expr("now()")).
		Set(columnUpdatedAt, sq.Expr("now()")).
		Where(sq.Eq{columnId: userId})

	query, err := db.BuildQuery("user_repo.delete", builder)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query)

	return err
}

func (r *Repository) ExistsById(ctx context.Context, userId int64) (bool, error) {
	builder := sq.
		Select("1").
		From(table).
		PlaceholderFormat(sq.Dollar).
		Prefix("SELECT EXISTS (").
		Where(sq.Eq{
			columnId:        userId,
			columnDeletedAt: nil,
		}).
		Suffix(")")

	query, err := db.BuildQuery("user_repo.exist_id", builder)
	if err != nil {
		return false, err
	}

	var isExists bool

	err = r.db.ScanOneContext(ctx, &isExists, query)
	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	builder := sq.
		Select("1").
		From(table).
		PlaceholderFormat(sq.Dollar).
		Prefix("SELECT EXISTS (").
		Where(sq.Eq{
			columnEmail:     email,
			columnDeletedAt: nil,
		}).
		Suffix(")")

	query, err := db.BuildQuery("user_repo.exist_email", builder)
	if err != nil {
		return false, err
	}

	var isExists bool

	err = r.db.ScanOneContext(ctx, &isExists, query)
	if err != nil {
		return false, err
	}

	return isExists, nil
}
