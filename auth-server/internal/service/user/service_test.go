package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/audit/dto"
	repoMocks "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/mocks"
	userDTO "github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestService_Get(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()

	userModel := model.FixtureUserWithFaker(t)

	userRepoMock := repoMocks.NewMockUser(t)
	userRepoMock.EXPECT().
		GetNotDeleted(ctx, userModel.ID).
		Return(userModel, nil)

	auditRepoMock := repoMocks.NewMockAudit(t)

	testTxManager := db.NewTestTxManager(t)
	service := New(log, userRepoMock, auditRepoMock, testTxManager)

	user, err := service.Get(ctx, userModel.ID)
	require.NoError(t, err)

	require.Equal(t, userModel, user)
}

func TestService_Create(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	userModel := model.FixtureUserWithFaker(t)

	userRepoMock := repoMocks.NewMockUser(t)
	userRepoMock.EXPECT().
		ExistsByEmail(ctx, userModel.Data.Email).
		Return(false, nil)

	userRepoMock.EXPECT().
		Create(ctx, mock.Anything).
		Return(userModel.ID, nil)

	auditRepoMock := repoMocks.NewMockAudit(t)
	auditRepoMock.EXPECT().
		Insert(ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, data dto.InsertDTO) (int64, error) {
			require.Equal(t, auditUserRef, data.Reference)
			require.Equal(t, userModel.ID, data.ReferenceID)
			require.Equal(t, auditNewUser, data.Action)

			return 1, nil
		})

	testTxManager := db.NewTestTxManager(t)
	service := New(log, userRepoMock, auditRepoMock, testTxManager)

	actualUserId, err := service.Create(ctx, userModel.Data)
	require.NoError(t, err)
	require.Equal(t, userModel.ID, actualUserId)
}

func TestService_Update(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	userModel := model.FixtureUserWithFaker(t)

	updateDTO := userDTO.UpdateDTO{
		Name:         sql.NullString{gofakeit.Name(), true},
		Email:        sql.NullString{gofakeit.Email(), true},
		PasswordHash: sql.NullString{gofakeit.HexUint16(), true},
		Role:         sql.NullInt32{gofakeit.Int32(), true},
	}

	userRepoMock := repoMocks.NewMockUser(t)
	userRepoMock.EXPECT().
		ExistsById(ctx, userModel.ID).
		Return(true, nil)

	userRepoMock.EXPECT().
		Update(ctx, userModel.ID, mock.Anything).
		Return(nil)

	auditRepoMock := repoMocks.NewMockAudit(t)
	auditRepoMock.EXPECT().
		Insert(ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, data dto.InsertDTO) (int64, error) {
			require.Equal(t, auditUserRef, data.Reference)
			require.Equal(t, userModel.ID, data.ReferenceID)
			require.Equal(t, auditUpdateUser, data.Action)

			return 1, nil
		})

	testTxManager := db.NewTestTxManager(t)
	service := New(log, userRepoMock, auditRepoMock, testTxManager)

	err = service.Update(ctx, userModel.ID, updateDTO)
	require.NoError(t, err)
}

func TestService_Delete(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	userModel := model.FixtureUserWithFaker(t)

	userRepoMock := repoMocks.NewMockUser(t)
	userRepoMock.EXPECT().
		ExistsById(ctx, userModel.ID).
		Return(true, nil)

	userRepoMock.EXPECT().
		Delete(ctx, userModel.ID).
		Return(nil)

	auditRepoMock := repoMocks.NewMockAudit(t)
	auditRepoMock.EXPECT().
		Insert(ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, data dto.InsertDTO) (int64, error) {
			require.Equal(t, auditUserRef, data.Reference)
			require.Equal(t, userModel.ID, data.ReferenceID)
			require.Equal(t, auditDeleteUser, data.Action)

			return 1, nil
		})

	testTxManager := db.NewTestTxManager(t)
	service := New(log, userRepoMock, auditRepoMock, testTxManager)

	err = service.Delete(ctx, userModel.ID)
	require.NoError(t, err)
}
