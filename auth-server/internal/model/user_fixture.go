package model

import (
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func FixtureUserWithFaker(t *testing.T) User {
	t.Helper()

	userModel := User{
		ID:        gofakeit.Int64(),
		CreatedAt: gofakeit.Date(),
		UpdatedAt: sql.NullTime{gofakeit.Date(), true},
		DeletedAt: sql.NullTime{},
		Data: UserData{
			Name:         gofakeit.Name(),
			Email:        gofakeit.Email(),
			PasswordHash: gofakeit.HexUint128(),
			Role:         gofakeit.Int32(),
		},
	}

	return userModel
}
