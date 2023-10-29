package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/config"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/db"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/logger"
	"github.com/aywan/balun_miserv_s2/chat-server/internal/runutil"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	desc.UnimplementedUserV1Server

	log *zap.Logger
	db  *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to open transaction", zap.Error(err))
		return nil, err
	}
	defer db.RollbackTx(ctx, tx, s.log, "error on rollback transaction")

	builderInsert := sq.Insert("\"chat\"").
		PlaceholderFormat(sq.Dollar).
		Columns("owner_id", "name").
		Values(req.OwnerId, req.Name).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	var chatID int64
	err = tx.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		s.log.Error("failed to insert chat", zap.Error(err))
		return nil, err
	}

	builderInsert = sq.Insert("\"message\"").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "type", "text").
		Values(chatID, desc.MessageType_SYSTEM, "start new chat").
		Suffix("RETURNING id")
	query, args, err = builderInsert.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}
	var msgID int64
	err = tx.QueryRow(ctx, query, args...).Scan(&msgID)
	if err != nil {
		s.log.Error("failed to insert message", zap.Error(err))
		return nil, err
	}

	builderInsert = sq.Insert("\"chat_message\"").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "last_message_id").
		Values(chatID, msgID)
	query, args, err = builderInsert.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error("failed to insert chat_message", zap.Error(err))
		return nil, err
	}

	builderInsert = sq.Insert("\"chat_user\"").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id", "last_message_id")
	for _, userID := range req.Users {
		builderInsert = builderInsert.Values(chatID, userID, msgID)
	}
	query, args, err = builderInsert.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error("failed to insert chat_user", zap.Error(err))
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return nil, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.ChatIdRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.
		Update("\"chat\"").
		PlaceholderFormat(sq.Dollar).
		Set("deleted_at", sq.Expr("now()")).
		Set("updated_at", sq.Expr("now()")).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error("failed to update chat", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to open transaction", zap.Error(err))
		return nil, err
	}
	defer db.RollbackTx(ctx, tx, s.log, "error on rollback transaction")

	builderInsert := sq.Insert("\"message\"").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id", "type", "text").
		Values(req.ChatId, req.UserId, req.Type, req.Text).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}
	var msgID int64
	err = tx.QueryRow(ctx, query, args...).Scan(&msgID)
	if err != nil {
		s.log.Error("failed to insert message", zap.Error(err))
		return nil, err
	}

	builderUpdate := sq.Update("\"chat_message\"").
		PlaceholderFormat(sq.Dollar).
		Set("last_message_id", sq.Expr("greatest(last_message_id, ?::bigint)", msgID)).
		Where(sq.Eq{"chat_id": req.ChatId})
	query, args, err = builderUpdate.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error("failed to update chat_message", zap.Error(err))
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetMessages(ctx context.Context, req *desc.GetMessagesRequest) (*desc.MessageListResponse, error) {
	builderSelect := sq.Select("id", "created_at", "user_id", "type", "text").
		From("\"message\"").
		PlaceholderFormat(sq.Dollar).
		Limit(uint64(req.Limit + 1))

	reverse := false

	if req.BeforeMessageId > 0 {
		builderSelect = builderSelect.
			Where(sq.Lt{"id": req.BeforeMessageId}).
			OrderBy("id DESC")
		reverse = true
	} else {
		builderSelect = builderSelect.
			Where(sq.Gt{"id": req.AfterMessageId}).
			OrderBy("id ASC")
	}

	query, args, err := builderSelect.ToSql()
	if err != nil {
		s.log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		s.log.Fatal("failed to select messages", zap.Error(err))
		return nil, err
	}

	items := make([]*desc.Message, 0)
	hasNext := false
	nextID := int64(0)

	var id int64
	var createdAt time.Time
	var userID sql.NullInt64
	var msgType desc.MessageType
	var text string

	cnt := int64(0)
	for rows.Next() {
		err = rows.Scan(&id, &createdAt, &userID, &msgType, &text)
		if err != nil {
			s.log.Error("failed to scan message", zap.Error(err))
			return nil, err
		}
		cnt++
		if cnt > req.Limit {
			hasNext = true
			nextID = id
			break
		}

		msg := new(desc.Message)
		msg.Id = id
		msg.CreatedAt = timestamppb.New(createdAt)
		if userID.Valid {
			msg.UserId = &userID.Int64
		}
		msg.Type = msgType
		msg.Text = text

		items = append(items, msg)
	}

	if reverse {
		slices.Reverse(items)
	}

	return &desc.MessageListResponse{
		Items:   items,
		HasNext: hasNext,
		NextId:  nextID,
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoadConfig()
	log := logger.MustNewLogger(cfg.Mode)

	dbPool, err := db.New(ctx, cfg.Db, log.With(zap.String("part", "db")))
	if err != nil {
		log.Fatal("failed to connect db", zap.Error(err))
	}

	serverPort, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	desc.RegisterUserV1Server(grpcServer, &server{
		db:  dbPool,
		log: log.With(zap.String("part", "server")),
	})
	go func() {
		if err := grpcServer.Serve(serverPort); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	time.Sleep(1 * time.Millisecond)

	err = testServer(ctx, serverPort.Addr(), log)
	if err != nil {
		log.Error("error test", zap.Error(err))
	}

	grpcServer.Stop()
}

func testServer(ctx context.Context, addr net.Addr, log *zap.Logger) error {
	conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer runutil.LogOnError(conn.Close, log, "error close client connection")

	c := desc.NewUserV1Client(conn)

	createReq, err := c.Create(ctx, &desc.CreateRequest{
		Users:   []int64{1, 2, 3},
		OwnerId: 1,
		Name:    "The chat",
	})
	if err != nil {
		return fmt.Errorf("failed to create chat, %w", err)
	}

	_, err = c.SendMessage(ctx, &desc.SendMessageRequest{
		ChatId: createReq.Id,
		UserId: 2,
		Text:   "text from two",
		Type:   desc.MessageType_USER,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	_, err = c.SendMessage(ctx, &desc.SendMessageRequest{
		ChatId: createReq.Id,
		UserId: 1,
		Text:   "text from one",
		Type:   desc.MessageType_USER,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	msgs, err := c.GetMessages(ctx, &desc.GetMessagesRequest{
		ChatId:          createReq.Id,
		Limit:           20,
		AfterMessageId:  0,
		BeforeMessageId: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}
	log.Info("get messages", zap.Any("messages", msgs.Items))

	_, err = c.Delete(ctx, &desc.ChatIdRequest{Id: createReq.Id})
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}
