package main

import (
	"context"
	"fmt"
	"net"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"homesearch.axel.to/commute/config"
	"homesearch.axel.to/commute/internal/clients"
	"homesearch.axel.to/commute/internal/service"
	"homesearch.axel.to/commute/internal/store"

	pb "homesearch.axel.to/services/commute/api"
	"homesearch.axel.to/shared/logger"
)

func main() {
	ctx := context.Background()
	log := logger.Init("commute")
	log.LogAttrs(ctx, logger.LevelInfo, "Starting commute service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	host := net.JoinHostPort("", cfg.Port)
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to listen on port %v: %v", cfg.Port, err))
		return
	}
	defer listen.Close()

	addressWrapperClient, err := clients.NewAddressWrapperServiceClient(ctx, cfg.AddressWrapperSvcHost)
	if err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to connect to address wrapper service at %v: %v", cfg.AddressWrapperSvcHost, err))
		return
	}
	defer addressWrapperClient.Close(ctx)

	roommateClient, err := clients.NewRoommateServiceClient(ctx, cfg.RoommateSvcHost)
	if err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to connect to roommate service at %v: %v", cfg.RoommateSvcHost, err))
		return
	}
	defer roommateClient.Close(ctx)

	commuteStore, err := store.NewCommuteStore(ctx, cfg)
	if err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to connect to commute store: %v", err))
		return
	}

	gRoutesServiceClient, err := routing.NewRoutesClient(ctx)
	if err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to connect to google routes service: %v", err))
		return
	}
	defer gRoutesServiceClient.Close()

	// TODO: initialize service with TLS connection
	service := service.NewCommuteService(addressWrapperClient, roommateClient, gRoutesServiceClient, commuteStore)

	grpcServer := grpc.NewServer()
	pb.RegisterCommuteServiceServer(grpcServer, service)

	log.LogAttrs(ctx, logger.LevelInfo, fmt.Sprintf("Starting server on port: %v", cfg.Port))

	if cfg.Environment == config.EnvironmentDevelopment {
		reflection.Register(grpcServer)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.LogAttrs(ctx, logger.LevelError, fmt.Sprintf("Failed to server request: %v", err))
	}
}
