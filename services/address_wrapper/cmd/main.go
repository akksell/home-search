package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"homesearch.axel.to/address_wrapper/config"
	"homesearch.axel.to/address_wrapper/internal/service"
	pb "homesearch.axel.to/services/address_wrapper/api"
)

func main() {
	// TODO: use proper logging tool -> see shared/go-pkg/logger
	//		 replace all usages of Printf with an implemented logger
	fmt.Printf("Starting Address Wrapper Service\n")

	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to start: %v\n", err)
	}
	address := net.JoinHostPort("", configuration.Port)

	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %v\n", configuration.Port)
	}
	defer listen.Close()

	// TODO: get credential file and create HTTPS connection
	service := &service.AddressWrapperService{
		Config: configuration,
	}
	// TODO: consider reading grpc options on initialization
	grpcServer := grpc.NewServer()
	pb.RegisterAddressWrapperServiceServer(grpcServer, service)
	fmt.Printf("Starting server on port: %v\n", configuration.Port)

	// TODO: disable this in prod via environment variables since
	//		 it opens up security vulnerabilities, its fine in local
	if configuration.Environment == config.EnvironmentDevelop {
		reflection.Register(grpcServer)
	}
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("Failed to server request: %v\n", err)
	}

}
