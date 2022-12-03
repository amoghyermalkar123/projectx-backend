package main

import (
	"context"
	"fmt"
	"projectx/proto/api"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	request := api.ChatMessageRequest{
		UserId:  1,
		Message: "Hello World!",
	}
	conn, err := grpc.Dial("127.0.0.1:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(2*time.Second),
	)

	if err != nil {
		fmt.Println(err)
	}

	chatServiceClient := api.NewChatMessageServiceClient(conn)
	chatS, err := chatServiceClient.Chat(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	err = chatS.Send(&request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Message sent successfully")
}
