package clients

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	base "homesearch.axel.to/base/types"
	pb "homesearch.axel.to/services/address_wrapper/api"
)

type AddressWrapperServiceClient struct {
	client pb.AddressWrapperServiceClient
	conn   *grpc.ClientConn
}

func NewAddressWrapperServiceClient(serviceURL string) (*AddressWrapperServiceClient, error) {
	var opts []grpc.DialOption

	// TODO: instantiate with TLS credentials
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serviceURL, opts...)
	if err != nil {
		return nil, err
	}

	c := &AddressWrapperServiceClient{
		client: pb.NewAddressWrapperServiceClient(conn),
		conn:   conn,
	}
	return c, nil
}

func (aws *AddressWrapperServiceClient) GetPlaceId(ctx context.Context, address *base.Address) (*pb.PlaceIdResponse, error) {
	requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	request := pb.PlaceIdRequest{
		Address: address,
	}

	return aws.client.GetPlaceId(requestCtx, &request)
}

// Close the grpc client connection
func (aws *AddressWrapperServiceClient) Close() {
	aws.conn.Close()
}
