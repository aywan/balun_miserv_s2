package repository

import (
	"context"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	auditDto "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/audit/dto"
	userDto "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
)

//go:generate mockery

type User interface {
	GetNotDeleted(ctx context.Context, userId int64) (model.User, error)
	Create(ctx context.Context, data model.UserData) (int64, error)
	Update(ctx context.Context, userId int64, data userDto.UpdateDTO) error
	Delete(ctx context.Context, userId int64) error
	ExistsById(ctx context.Context, userId int64) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type Audit interface {
	Insert(ctx context.Context, data auditDto.InsertDTO) (int64, error)
}
