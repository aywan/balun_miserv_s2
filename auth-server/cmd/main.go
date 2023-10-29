package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/config"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/db"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/logger"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/runutil"
	"github.com/aywan/balun_miserv_s2/auth-server/internal/security"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	sq "github.com/Masterminds/squirrel"
	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	desc.UnimplementedUserV1Server
	log    *zap.Logger
	dbPool *pgxpool.Pool
}

func (s *server) Get(ctx context.Context, req *desc.UserIdRequest) (*desc.UserResponse, error) {
	builderSelect := sq.
		Select("id", "created_at", "updated_at", "name", "email", "role").
		From("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"id":         req.Id,
			"deleted_at": nil,
		})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	rsp := &desc.UserResponse{User: &desc.UserData{}}
	var createdAt, updatedAt sql.NullTime

	err = s.dbPool.
		QueryRow(ctx, query, args...).
		Scan(&rsp.Id, &createdAt, &updatedAt, &rsp.User.Name, &rsp.User.Email, &rsp.User.Role)
	if err != nil {
		s.log.Error("failed to select user", zap.Error(err))
		return nil, err
	}

	if createdAt.Valid {
		rsp.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		rsp.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return rsp, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.Credentials.Password != req.Credentials.PasswordConfirm {
		return nil, errors.New("password and confirm not equal")
	}

	passwordHash, err := security.HashPassword(req.Credentials.Password)
	if err != nil {
		return nil, err
	}

	builderInsert := sq.Insert("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role", "password_hash").
		Values(req.User.Name, req.User.Email, req.User.Role, passwordHash).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	var userID int64
	err = s.dbPool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		s.log.Error("failed to insert user", zap.Error(err))
		return nil, err
	}

	return &desc.CreateResponse{Id: userID}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.
		Update("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", sq.Expr("now()")).
		Where(sq.Eq{"id": req.Id})

	if req.Email != nil {
		builderUpdate = builderUpdate.Set("email", req.Email.Value)
	}
	if req.Name != nil {
		builderUpdate = builderUpdate.Set("name", req.Name.Value)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	_, err = s.dbPool.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error("failed to update user", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.UserIdRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.
		Update("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Set("deleted_at", sq.Expr("now()")).
		Set("updated_at", sq.Expr("now()")).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	_, err = s.dbPool.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error("failed to delete user", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func main() {
	cfg := config.MustLoadConfig()
	log := logger.MustNewLogger(cfg.Mode)
	defer runutil.IgnoreErr(log.Sync)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dbPool, err := db.New(
		ctx,
		cfg.Db,
		log.With(zap.String("part", "db")),
	)
	if err != nil {
		log.Fatal("failed to connect db", zap.Error(err))
	}
	defer dbPool.Close()

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatal("failed to ping db", zap.Error(err))
	}

	serverPort, err := net.Listen("tcp", cfg.Server.Listen)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	desc.RegisterUserV1Server(
		grpcServer,
		&server{
			log:    log.With(zap.String("part", "server")),
			dbPool: dbPool,
		},
	)
	go func() {
		if err = grpcServer.Serve(serverPort); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	time.Sleep(1 * time.Millisecond)

	err = testConnectToServer(ctx, serverPort.Addr().String(), log)
	if err != nil {
		log.Error("failed to test", zap.Error(err))
	}

	grpcServer.Stop()
}

func testConnectToServer(ctx context.Context, serverAddr string, log *zap.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer runutil.LogOnError(conn.Close, log, "close client connection")

	c := desc.NewUserV1Client(conn)

	password := security.CreatePassword(6)
	rndEmail := fmt.Sprintf("some+%d@somemail.ru", rand.Int()) // #nosec G404 -- allow.
	createRsp, err := c.Create(ctx, &desc.CreateRequest{
		User: &desc.UserData{
			Name:  "SomeOne",
			Email: rndEmail,
			Role:  desc.UserRole_USER,
		},
		Credentials: &desc.UserCredentials{
			Password:        password,
			PasswordConfirm: password,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	getRsp, err := c.Get(ctx, &desc.UserIdRequest{Id: createRsp.Id})
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	log.Info(color.GreenString("id: %d", getRsp.GetId()))
	log.Info(color.GreenString("data: %+v", getRsp))

	rndEmail = fmt.Sprintf("new-some+%d@somemail.ru", rand.Int()) // #nosec G404 -- allow.
	_, err = c.Update(ctx, &desc.UpdateRequest{
		Id:    createRsp.Id,
		Email: wrapperspb.String(rndEmail),
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	_, err = c.Delete(ctx, &desc.UserIdRequest{Id: createRsp.Id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	_, err = c.Get(ctx, &desc.UserIdRequest{Id: createRsp.Id})
	if err == nil {
		return fmt.Errorf("user found after delete")
	}

	return nil
}
