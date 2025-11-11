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

	pb "homesearch.axel.to/services/roommate/api"
	"homesearch.axel.to/shared/logger"
)

var roomateTokenSource oauth2.TokenSource

type RoommateServiceClient struct {
	client pb.RoommateServiceClient
	conn   *grpc.ClientConn
}

func NewRoommateServiceClient(ctx context.Context, serviceURL string) (*RoommateServiceClient, error) {
	var opts []grpc.DialOption
	host := strings.Split(serviceURL, ":")[0]
	port := strings.Split(serviceURL, ":")[1]
	isSecure := port == "443"

	if !isSecure {
		logger.LogAttrs(ctx, logger.LevelInfo, "roommate service connection is insecure")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		logger.LogAttrs(ctx, logger.LevelInfo, "securely connecting to roommate service with TLS")
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

	if roomateTokenSource == nil {
		var protocol string = "http://"
		if isSecure {
			protocol = "https://"
		}
		roomateTokenSource, err = idtoken.NewTokenSource(ctx, protocol+host)
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
		logger.LogAttrs(ctx, logger.LevelError, "failed to get token from roomate token source")
		return nil, fmt.Errorf("Failed to get group points of interest: %w", err)
	}

	reqContext = grpcMetadata.AppendToOutgoingContext(reqContext, "authorization", "Bearer "+token.AccessToken)

	request := &pb.GetGroupPOIsRequest{
		GroupId: groupId,
	}

	logger.LogAttrs(reqContext, logger.LevelInfo, "requesting points of interest", logger.String("groupId", groupId))
	return rs.client.GetGroupPointsOfInterest(reqContext, request)
}

func (rs *RoommateServiceClient) Close(ctx context.Context) {
	logger.LogAttrs(ctx, logger.LevelInfo, "closing roommate service connection")
	rs.conn.Close()
}
