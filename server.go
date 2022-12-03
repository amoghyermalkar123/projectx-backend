package main

import (
	"fmt"
	"log"
	"net"
	"projectx/proto/api"

	"google.golang.org/grpc"
)

type MessageServer struct{}

func (m *MessageServer) Chat(server api.ChatMessageService_ChatServer) error {
	for {
		request, err := server.Recv()
		if err != nil {
			fmt.Println("err from chat", err)
			return err
		}
		err = server.Send(&api.ChatMessageResponse{
			UserId:  request.UserId,
			Message: fmt.Sprint("Response from Keyur:", request.Message),
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func NewRpcServer() *MessageServer {
	return &MessageServer{}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8081))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterChatMessageServiceServer(s, NewRpcServer())
	log.Printf("Starting gRPC server on tcp port %d\n", 8080)
	log.Fatal(s.Serve(lis))
}
