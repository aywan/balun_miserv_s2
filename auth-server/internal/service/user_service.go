package service

import (
	"context"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
)

type User interface {
	Get(ctx context.Context, userId int64) (model.User, error)
	Create(ctx context.Context, data model.UserData) (int64, error)
	Update(ctx context.Context, userId int64, data dto.UpdateDTO) error
	Delete(ctx context.Context, userId int64) error
}
