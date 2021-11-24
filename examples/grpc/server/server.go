package main

import (
	"log"
	"net"
	"time"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterClockServiceServer(grpcServer, new(ClockService))
	log.Fatal(grpcServer.Serve(lis))
}

type ClockService struct {
	UnimplementedClockServiceServer
}

func (s *ClockService) GetTime(ctx context.Context, req *GetTimeRequest) (*GetTimeResponse, error) {
	return &GetTimeResponse{FormattedTime: time.Now().String()}, nil
}
