package clients

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"

	base "homesearch.axel.to/base/types"
	pb "homesearch.axel.to/services/address_wrapper/api"
	"homesearch.axel.to/shared/logger"
)

var addressTokenSource oauth2.TokenSource

type AddressWrapperServiceClient struct {
	client pb.AddressWrapperServiceClient
	conn   *grpc.ClientConn
}

func NewAddressWrapperServiceClient(ctx context.Context, serviceURL string) (*AddressWrapperServiceClient, error) {
	var opts []grpc.DialOption
	host := strings.Split(serviceURL, ":")[0]
	port := strings.Split(serviceURL, ":")[1]
	isSecure := port == "443"

	if !isSecure {
		logger.LogAttrs(ctx, logger.LevelInfo, "address wrapper service connection is insecure")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		logger.LogAttrs(ctx, logger.LevelInfo, "securely connecting to address wrapper service with TLS")
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	conn, err := grpc.NewClient(serviceURL, opts...)
	if err != nil {
		return nil, err
	}

	if addressTokenSource == nil {
		var protocol string = "http://"
		if isSecure {
			protocol = "https://"
		}
		addressTokenSource, err = idtoken.NewTokenSource(ctx, protocol+host)
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

	logger.LogAttrs(requestCtx, logger.LevelInfo, "get placeId from address", logger.Group("address", address))
	return aws.client.GetPlaceId(requestCtx, &request)
}

// Close the grpc client connection
func (aws *AddressWrapperServiceClient) Close(ctx context.Context) {
	logger.LogAttrs(ctx, logger.LevelInfo, "closing address wrapper service connection")
	aws.conn.Close()
}
