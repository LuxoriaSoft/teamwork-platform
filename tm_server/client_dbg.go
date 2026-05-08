package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "server/msgpb"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:51234",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	stream, err := client.MsgStream(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	// Receive messages
	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("[%s]: %s\n", msg.User, msg.Message)
		}
	}()

	// Send messages
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()

		err := stream.Send(&pb.MsgMessage{
			User:    "Client",
			Message: text,
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}