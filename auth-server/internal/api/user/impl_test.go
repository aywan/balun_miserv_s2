package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/service/mocks"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestImplementation_Create(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	userID := gofakeit.Int64()
	pwd := gofakeit.HexUint16()

	req := &desc.CreateRequest{
		User: &desc.UserData{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Role:  desc.UserRole_ADMIN,
		},
		Credentials: &desc.UserCredentials{
			Password:        pwd,
			PasswordConfirm: pwd,
		},
	}

	userService := mocks.NewMockUser(t)
	userService.EXPECT().
		Create(ctx, model.UserData{
			Name:         req.User.Name,
			Email:        req.User.Email,
			PasswordHash: req.Credentials.Password,
			Role:         int32(req.User.Role),
		}).
		Return(userID, nil)

	srv := New(log, userService)

	rsp, err := srv.Create(ctx, req)
	require.NoError(t, err)

	require.Equal(t, rsp.Id, userID)
}

func TestImplementation_Update(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	req := desc.UpdateRequest{
		Id:    gofakeit.Int64(),
		Name:  wrapperspb.String(gofakeit.Name()),
		Email: wrapperspb.String(gofakeit.Email()),
	}

	userService := mocks.NewMockUser(t)
	userService.EXPECT().
		Update(ctx, req.Id, dto.UpdateDTO{
			Name:         sql.NullString{req.Name.Value, true},
			Email:        sql.NullString{req.Email.Value, true},
			PasswordHash: sql.NullString{},
			Role:         sql.NullInt32{},
		}).
		Return(nil)

	srv := New(log, userService)

	_, err = srv.Update(ctx, &req)
	require.NoError(t, err)
}

func TestImplementation_Delete(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	userID := gofakeit.Int64()

	userService := mocks.NewMockUser(t)
	userService.EXPECT().
		Delete(ctx, userID).
		Return(nil)

	srv := New(log, userService)

	_, err = srv.Delete(ctx, &desc.UserIdRequest{Id: userID})
	require.NoError(t, err)
}

func TestImplementation_Get(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	var err error

	userModel := model.FixtureUserWithFaker(t)

	userService := mocks.NewMockUser(t)
	userService.EXPECT().
		Get(ctx, userModel.ID).
		Return(userModel, nil)

	srv := New(log, userService)

	rsp, err := srv.Get(ctx, &desc.UserIdRequest{Id: userModel.ID})
	require.NoError(t, err)

	require.Equal(t, rsp.Id, userModel.ID)
	require.Equal(t, rsp.CreatedAt, timestamppb.New(userModel.CreatedAt))
	require.Equal(t, rsp.UpdatedAt, timestamppb.New(userModel.UpdatedAt.Time))
	require.Equal(t, rsp.User.Name, userModel.Data.Name)
	require.Equal(t, rsp.User.Email, userModel.Data.Email)
	require.Equal(t, rsp.User.Role, desc.UserRole(userModel.Data.Role))
}
