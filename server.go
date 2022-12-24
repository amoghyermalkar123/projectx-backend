package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"projectx/proto/api"

	"google.golang.org/grpc"
)

type MessageServer struct {
	channel map[string][]chan *api.ChatMessage
}

func (m *MessageServer) JoinChannel(joinRequest *api.ConnectionRequest, stream api.ChatMessageService_JoinChannelServer) error {
	msgChannel := make(chan *api.ChatMessage)
	m.channel[joinRequest.ChannelName] = append(m.channel[joinRequest.ChannelName], msgChannel)
	log.Println("channel data:", m.channel)
	for {
		select {
		case <-stream.Context().Done():
			log.Println("closing stream")
			return nil
		case msg := <-msgChannel:
			log.Println("streaming msg to:", joinRequest.UserID)
			err := stream.Send(msg)
			if err != nil {
				log.Println("error while streaming response to client", err)
			}
		}
	}
}

func (m *MessageServer) SendMessage(server api.ChatMessageService_SendMessageServer) error {
	for {
		msg, err := server.Recv()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			if err == io.EOF {
				return nil
			}
			log.Println("error while receing client message", err)
			return err
		}

		go func() {
			streams := m.channel[msg.ChannelName]
			for _, individualUserStream := range streams {
				individualUserStream <- msg
			}
		}()
	}
}

func NewRpcServer() *MessageServer {
	return &MessageServer{
		channel: make(map[string][]chan *api.ChatMessage),
	}
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
