package main

import (
	"log"
	"net"

	"github.com/moondoggy/courier/internal/service"
	courierv1 "github.com/moondoggy/courier/proto/courier/v1"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()
	courierv1.RegisterCourierServiceServer(s, service.NewServer())

	log.Println("listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
