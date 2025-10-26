package clients

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "homesearch.axel.to/services/roommate/api"
)

type RoommateServiceClient struct {
	client pb.RoommateServiceClient
	conn   *grpc.ClientConn
}

func NewRoommateServiceClient(serviceURL string) (*RoommateServiceClient, error) {
	var opts []grpc.DialOption

	// TODO: instantiate with TLS credentials
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serviceURL, opts...)
	if err != nil {
		return nil, err
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

	request := &pb.GetGroupPOIsRequest{
		GroupId: groupId,
	}

	return rs.client.GetGroupPointsOfInterest(reqContext, request)
}

func (rs *RoommateServiceClient) Close() {
	rs.conn.Close()
}
