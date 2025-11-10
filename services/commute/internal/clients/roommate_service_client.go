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

	pb "homesearch.axel.to/services/roommate/api"
)

var roomateTokenSource oauth2.TokenSource

type RoommateServiceClient struct {
	client pb.RoommateServiceClient
	conn   *grpc.ClientConn
}

func NewRoommateServiceClient(ctx context.Context, serviceURL string) (*RoommateServiceClient, error) {
	var opts []grpc.DialOption

	// TODO: instantiate with TLS credentials
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serviceURL, opts...)
	if err != nil {
		return nil, err
	}

	if roomateTokenSource == nil {
		roomateTokenSource, err = idtoken.NewTokenSource(ctx, serviceURL)
		if err != nil {
			return nil, fmt.Errorf("Failed to start roommate service client: %w", err)
		}
	}

	c := &RoommateServiceClient{
		client: pb.NewRoommateServiceClient(conn),
		conn:   conn,
	}
	return c, nil
}

func (rs *RoommateServiceClient) GetGroupPointsOfInterest(ctx context.Context, groupId string) (*pb.GetGroupPOIsResponse, error) {
	reqContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	token, err := roomateTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("Failed to get group points of interest: %w", err)
	}

	reqContext = grpcMetadata.AppendToOutgoingContext(reqContext, "audience", "Bearer "+token.AccessToken)

	request := &pb.GetGroupPOIsRequest{
		GroupId: groupId,
	}

	return rs.client.GetGroupPointsOfInterest(reqContext, request)
}

func (rs *RoommateServiceClient) Close() {
	rs.conn.Close()
}
