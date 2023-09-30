package main

import (
	"context"
	"log"
	"net"
	"time"

	desc "github.com/aywan/balun_miserv_s2/auth-server/pkg/grpc/v1/user_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Get(_ context.Context, req *desc.UserIdRequest) (*desc.UserResponse, error) {
	log.Println(color.RedString("request user id: %d", req.GetId()))

	return &desc.UserResponse{
		Id:        123,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
		User: &desc.UserData{
			Name:  "a",
			Email: "b",
			Role:  desc.UserRole_USER,
		},
	}, nil
}

func main() {
	serverPort, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	desc.RegisterUserV1Server(grpcServer, &server{})
	go func() {
		if err = grpcServer.Serve(serverPort); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(1 * time.Millisecond)

	conn, err := grpc.Dial(serverPort.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("err: %v", err)
		}
	}()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.UserIdRequest{Id: 123})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Println(color.GreenString("id: %d", r.GetId()))
	log.Println(color.GreenString("data: %+v", r.GetUser()))

	grpcServer.Stop()
}
