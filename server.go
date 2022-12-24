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
	channel map[string]map[string]api.ChatMessageService_JoinChannelServer
}

func (m *MessageServer) JoinChannel(joinRequest *api.ConnectionRequest, stream api.ChatMessageService_JoinChannelServer) error {
	if _, ok := m.channel[joinRequest.ChannelName]; !ok {
		m.channel[joinRequest.ChannelName] = map[string]api.ChatMessageService_JoinChannelServer{
			joinRequest.UserID: stream,
		}
	} else {
		m.channel[joinRequest.ChannelName][joinRequest.UserID] = stream
	}
	fmt.Println(m.channel)
	for {
		select {
		case <-stream.Context().Done():
			log.Println("closing stream")
			m.channel[joinRequest.ChannelName][joinRequest.UserID] = nil
			return nil
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
			for _, stream := range streams {
				if stream == nil {
					log.Println("client is disconnected, not sending message")
					return
				}
				err := stream.Send(msg)
				if err != nil {
					log.Println("error while sending message to channel", err)
					return
				}
			}
		}()
	}
}

func NewRpcServer() *MessageServer {
	m := make(map[string]map[string]api.ChatMessageService_JoinChannelServer)
	return &MessageServer{
		channel: m,
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
