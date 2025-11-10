package clients

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"

	base "homesearch.axel.to/base/types"
	pb "homesearch.axel.to/services/address_wrapper/api"
)

var addressTokenSource oauth2.TokenSource

type AddressWrapperServiceClient struct {
	client pb.AddressWrapperServiceClient
	conn   *grpc.ClientConn
}

func NewAddressWrapperServiceClient(ctx context.Context, serviceURL string) (*AddressWrapperServiceClient, error) {
	var opts []grpc.DialOption

	// TODO: instantiate with TLS credentials
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serviceURL, opts...)
	if err != nil {
		return nil, err
	}

	if addressTokenSource == nil {
		addressTokenSource, err = idtoken.NewTokenSource(ctx, serviceURL)
		if err != nil {
			return nil, fmt.Errorf("Failed to start address wrapper client: %w", err)
		}
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

	token, err := addressTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("Failed to get place id: %w", err)
	}

	requestCtx = grpcMetadata.AppendToOutgoingContext(requestCtx, "authorization", "Bearer "+token.AccessToken)

	request := pb.PlaceIdRequest{
		Address: address,
	}

	return aws.client.GetPlaceId(requestCtx, &request)
}

// Close the grpc client connection
func (aws *AddressWrapperServiceClient) Close() {
	aws.conn.Close()
}
