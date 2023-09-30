package main

import (
	"context"
	"log"
	"net"
	"time"

	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Println(color.RedString("request usernames: %v", req.GetUsernames()))

	return &desc.CreateResponse{
		Id: 123,
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

	r, err := c.Create(ctx, &desc.CreateRequest{Usernames: []string{"x", "y"}})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Println(color.GreenString("new chat id: %d", r.GetId()))

	grpcServer.Stop()
}
