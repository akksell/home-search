package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"homesearch.axel.to/roommate/config"
	"homesearch.axel.to/roommate/internal/clients"
	"homesearch.axel.to/roommate/internal/service"
	"homesearch.axel.to/roommate/internal/store"
	pb "homesearch.axel.to/services/roommate/api"
)

func main() {
	fmt.Printf("Starting Roommate Service...\n")

	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to start Roommate Service: %v\n", err)
		return
	}

	host := net.JoinHostPort("", cfg.Port)
	listen, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Printf("Failed to start Roommate Service: %v\n", err)
		return
	}
	defer listen.Close()

	roommateStore, err := store.NewRoommateStore(ctx, cfg)
	if err != nil {
		fmt.Printf("Failed to start Roommate Service: %v\n", err)
		return
	}
	defer roommateStore.Close()

	addressWrapperSvcClient, err := clients.NewAddressWrapperServiceClient(cfg.AddressWrapperServiceHost)
	if err != nil {
		fmt.Printf("Failed to start Roommate Service: %v\n", err)
		return
	}
	defer addressWrapperSvcClient.Close()

	opts := make([]grpc.ServerOption, 0)

	server := service.NewRoomateService(roommateStore, addressWrapperSvcClient)
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRoommateServiceServer(grpcServer, server)

	fmt.Printf("Starting server on port: %v\n", cfg.Port)
	// TODO: disable this in prod via environment variables since
	//		 it opens up security vulnerabilities, its fine in local
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("Failed to server request: %v", err)
	}
}
