package converter

import (
	"database/sql"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/model"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/repository/user/dto"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToGrpcUserResponse(user model.User) *desc.UserResponse {
	rsp := &desc.UserResponse{User: &desc.UserData{}}

	rsp.Id = user.ID
	rsp.User.Name = user.Data.Name
	rsp.User.Email = user.Data.Email
	rsp.User.Role = desc.UserRole(user.Data.Role)
	rsp.CreatedAt = timestamppb.New(user.CreatedAt)
	if user.UpdatedAt.Valid {
		rsp.UpdatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return rsp
}

func GrpcCreateRequestToUserData(req *desc.CreateRequest) model.UserData {
	return model.UserData{
		Name:         req.User.Name,
		Email:        req.User.Email,
		PasswordHash: req.Credentials.Password,
		Role:         int32(req.User.Role),
	}
}

func GrpcUpdateRequestToUpdateDTO(req *desc.UpdateRequest) dto.UpdateDTO {
	out := dto.UpdateDTO{
		Name:         sql.NullString{},
		Email:        sql.NullString{},
		PasswordHash: sql.NullString{},
		Role:         sql.NullInt32{},
	}

	if req.Name != nil {
		out.Name.String = req.Name.Value
		out.Name.Valid = true
	}

	if req.Email != nil {
		out.Email.String = req.Email.Value
		out.Email.Valid = true
	}

	return out
}
