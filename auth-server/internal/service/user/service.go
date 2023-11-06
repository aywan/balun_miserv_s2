package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository"
	dto2 "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/audit/dto"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/security"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/service"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"go.uber.org/zap"
)

const auditUserRef = "user"

type Service struct {
	log       *zap.Logger
	userRepo  repository.User
	auditRepo repository.Audit
	txManager db.TxManager
}

var _ service.User = (*Service)(nil)

func New(
	logger *zap.Logger,
	userRepo repository.User,
	auditRepo repository.Audit,
	txManager db.TxManager,
) *Service {
	return &Service{
		log:       logger,
		userRepo:  userRepo,
		auditRepo: auditRepo,
		txManager: txManager,
	}
}

func (s *Service) Get(ctx context.Context, userId int64) (model.User, error) {
	return s.userRepo.GetNotDeleted(ctx, userId)
}

func (s *Service) Create(ctx context.Context, data model.UserData) (int64, error) {
	passwordHash, err := security.HashPassword(data.PasswordHash)
	if err != nil {
		return 0, err
	}
	data.PasswordHash = passwordHash

	var retUserId int64

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		isExist, err := s.userRepo.ExistsByEmail(ctx, data.Email)
		if err != nil {
			return err
		}
		if isExist {
			return fmt.Errorf("user with email %s already exists", data.Email)
		}

		userId, err := s.userRepo.Create(ctx, data)
		if err != nil {
			return err
		}

		_, err = s.auditRepo.Insert(ctx, dto2.InsertDTO{
			CreatorId:   sql.NullInt64{},
			Reference:   auditUserRef,
			ReferenceID: userId,
			Action:      "new user",
		})

		if err != nil {
			return err
		}

		retUserId = userId
		return nil
	})

	return retUserId, err
}

func (s *Service) Update(ctx context.Context, userId int64, data dto.UpdateDTO) error {
	isExist, err := s.userRepo.ExistsById(ctx, userId)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("user with id %d does not exist", userId)
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepo.Update(ctx, userId, data)
		if err != nil {
			return err
		}

		_, err = s.auditRepo.Insert(ctx, dto2.InsertDTO{
			CreatorId:   sql.NullInt64{},
			Reference:   auditUserRef,
			ReferenceID: userId,
			Action:      "update user",
		})

		return err
	})
}

func (s *Service) Delete(ctx context.Context, userId int64) error {
	isExist, err := s.userRepo.ExistsById(ctx, userId)
	if err != nil {
		return err
	}
	if !isExist {
		return nil
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepo.Delete(ctx, userId)
		if err != nil {
			return err
		}

		_, err = s.auditRepo.Insert(ctx, dto2.InsertDTO{
			CreatorId:   sql.NullInt64{},
			Reference:   auditUserRef,
			ReferenceID: userId,
			Action:      "delete user",
		})

		return err
	})
}
