package main

import (
	"context"
	"fmt"
	"log"
	"net"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"homesearch.axel.to/commute/config"
	"homesearch.axel.to/commute/internal/clients"
	"homesearch.axel.to/commute/internal/service"
	"homesearch.axel.to/commute/internal/store"

	pb "homesearch.axel.to/services/commute/api"
)

func main() {
	fmt.Println("Starting commute service...")

	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("%w", err)
		return
	}
	listen, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		fmt.Printf("Failed to listen on port %v: %v\n", cfg.Port, err)
		return
	}
	defer listen.Close()

	addressWrapperClient, err := clients.NewAddressWrapperServiceClient(cfg.AddressWrapperSvcHost)
	if err != nil {
		fmt.Printf("Failed to connect to address wrapper service at %v: %v\n", cfg.AddressWrapperSvcHost, err)
		return
	}
	defer addressWrapperClient.Close()

	roommateClient, err := clients.NewRoommateServiceClient(cfg.RoommateSvcHost)
	if err != nil {
		fmt.Printf("Failed to connect to roommate service at %v: %v\n", cfg.RoommateSvcHost, err)
		return
	}
	defer roommateClient.Close()

	commuteStore, err := store.NewCommuteStore(ctx, cfg)
	if err != nil {
		fmt.Println("Failed to connect to commute store: %w", err)
		return
	}

	gRoutesServiceClient, err := routing.NewRoutesClient(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to google routes service: %w\n", err)
		return
	}
	defer gRoutesServiceClient.Close()

	// TODO: initialize service with TLS connection
	service := service.NewCommuteService(addressWrapperClient, roommateClient, gRoutesServiceClient, commuteStore)

	grpcServer := grpc.NewServer()
	pb.RegisterCommuteServiceServer(grpcServer, service)

	fmt.Printf("Starting server on port: %v\n", cfg.Port)

	// TODO: disable this in prod via environment variables since
	//		 it opens up security vulnerabilities, its fine in local
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("Failed to server request: %v\n", err)
	}
}
