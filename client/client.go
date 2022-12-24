package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"projectx/proto/api"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(2*time.Second),
	)

	if err != nil {
		fmt.Println(err)
	}

	chatServiceClient := api.NewChatMessageServiceClient(conn)

	channelStream, err := chatServiceClient.JoinChannel(context.Background(), &api.ConnectionRequest{
		UserID:      "User2",
		ChannelName: "Amogh",
	})
	fmt.Println("joined channel")
	if err != nil {
		log.Println("error while joining to channel", err)
		return
	}

	go func() {
		fmt.Println("started listening for message")
		for {
			fmt.Println("waiting for message ..")
			msg, err := channelStream.Recv()
			if err == io.EOF {
				log.Print("end")
				return
			}
			log.Println("message received for user from channel:", msg)
		}
	}()
	request := api.ChatMessage{
		ChannelName: "Amogh",
		Message:     "hi there",
	}
	stream, err := chatServiceClient.SendMessage(context.Background())

	if err != nil {
		log.Println("error while joining to channel", err)
		return
	}

	for {
		time.Sleep(2 * time.Second)
		err = stream.Send(&request)
		if err != nil {
			log.Println("error while sending message", err)
			return
		}
	}
}
