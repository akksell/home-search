package main

import (
	"fmt"
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
	fmt.Printf("Starting Address Wrapper Service")

	config := config.LoadConfig()
	// TODO: read this from env/cloud run deployment
	port := ":8080"

	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Failed to listen on port %v", port)
	}
	defer listen.Close()

	// TODO: get credential file and create HTTPS connection
	service := &service.AddressWrapperService{
		Config: config,
	}
	// TODO: consider reading grpc options on initialization
	grpcServer := grpc.NewServer()
	pb.RegisterAddressWrapperServiceServer(grpcServer, service)
	fmt.Printf("Starting server on port: %v", port)

	// TODO: disable this in prod via environment variables since
	//		 it opens up security vulnerabilities, its fine in local
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("Failed to server request: %v", err)
	}

}
