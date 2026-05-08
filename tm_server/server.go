package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "server/msgpb"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

func (s *server) MsgStream(stream pb.MessageService_MsgStreamServer) error {
	log.Println("Client connected")

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			log.Println("Client disconnected")
			return nil
		}

		if err != nil {
			return err
		}

		log.Printf("[%s]: %s\n", msg.User, msg.Message)

		response := &pb.MsgMessage{
			User:    "Server",
			Message: "Received: " + msg.Message,
		}

		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":51234")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterMessageServiceServer(grpcServer, &server{})

	log.Println("Server running on :51234")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}