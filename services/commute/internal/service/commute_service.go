package service

import (
	"context"

	pb "homesearch.axel.to/services/commute/api"
)

type commuteService struct {
	pb.UnimplementedCommuteServiceServer
}

func (s *commuteService) Calculate(ctx context.Context, request *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	response := &pb.CalculateResponse{}
	return response, nil
}
