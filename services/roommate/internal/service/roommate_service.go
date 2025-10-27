package service

import (
	"context"

	"homesearch.axel.to/roommate/internal/clients"
	"homesearch.axel.to/roommate/internal/store"
	"homesearch.axel.to/roommate/internal/store/models"
	pb "homesearch.axel.to/services/roommate/api"
)

type RoommateService struct {
	roommateStore     *store.RoommateStore
	addressWrapperSvc *clients.AddressWrapperServiceClient
	pb.UnimplementedRoommateServiceServer
}

func NewRoomateService(roommateStore *store.RoommateStore, addressWrapperSvcClient *clients.AddressWrapperServiceClient) *RoommateService {
	return &RoommateService{
		roommateStore:     roommateStore,
		addressWrapperSvc: addressWrapperSvcClient,
	}
}

func (s *RoommateService) AddRoommate(ctx context.Context, req *pb.AddRoommateRequest) (*pb.AddRoommateResponse, error) {
	// TODO: add validation
	roommate := req.GetRoommate()
	person := roommate.GetPersonalDetails()
	roommateModel := &models.Roommate{
		FirstName: person.GetFirstName(),
		LastName:  person.GetLastName(),
	}
	id, err := s.roommateStore.CreateRoommate(ctx, roommateModel)
	if err != nil {
		// TODO: return error message here
		return &pb.AddRoommateResponse{
			RoommateId: "",
		}, err
	}
	return &pb.AddRoommateResponse{
		RoommateId: id,
	}, nil
}

func (s *RoommateService) UpdateRoommate(ctx context.Context, req *pb.UpdateRoommateRequest) (*pb.UpdateRoommateResponse, error) {
	return &pb.UpdateRoommateResponse{}, nil
}

func (s *RoommateService) RetireRoommate(ctx context.Context, req *pb.RetireRoommateRequest) (*pb.RetireRoommateResponse, error) {
	return &pb.RetireRoommateResponse{}, nil
}

func (s *RoommateService) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	// TODO: validation
	// TODO: get ownerId from the request headers
	group := &models.RoommateGroup{
		Name:        req.GetGroupName(),
		Description: req.GetGroupDescription(),
	}
	groupId, err := s.roommateStore.CreateGroup(ctx, group)
	if err != nil {
		return &pb.CreateGroupResponse{
			JoinUrl: "",
		}, err
	}
	// TODO: generate a new join url. For now, just return the groupId
	return &pb.CreateGroupResponse{
		JoinUrl: groupId,
	}, nil
}

func (s *RoommateService) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	return &pb.JoinGroupResponse{
		Status: "TODO",
	}, nil
}

func (s *RoommateService) CreatePointOfInterest(ctx context.Context, req *pb.CreatePointOfInterestRequest) (*pb.CreatePointOfInterestResponse, error) {
	// TODO: validation
	placeId, err := s.addressWrapperSvc.GetPlaceId(ctx, req.GetAddress())
	if err != nil {
		// TODO: error response
		return &pb.CreatePointOfInterestResponse{}, err
	}

	pointOfInterestAddress := &models.Address{
		Street:            req.GetAddress().GetStreet(),
		City:              req.GetAddress().GetCity(),
		PostalCode:        req.GetAddress().GetPostalCode(),
		StateProvinceCode: req.GetAddress().GetStateProvinceCode(),
		CountryCode:       req.GetAddress().GetCountryCode(),
	}
	pointOfInterest := &models.PointOfInterest{
		Label:   req.GetLabel(),
		Tags:    req.GetTags(),
		PlaceId: placeId.GetPlaceId(),
		Address: pointOfInterestAddress,
	}
	id, err := s.roommateStore.CreatePointOfInterest(ctx, pointOfInterest)
	if err != nil {
		// TODO: add error message in response
		return &pb.CreatePointOfInterestResponse{}, err
	}

	return &pb.CreatePointOfInterestResponse{
		Id:      id,
		PlaceId: placeId.GetPlaceId(),
	}, nil
}

func (s *RoommateService) AddPointOfInterest(ctx context.Context, req *pb.AddPointOfInterestRequest) (*pb.AddPointOfInterestResponse, error) {
	// TODO: validation
	pointOfInterest, err := s.roommateStore.GetPointOfInterest(ctx, req.GetPointOfInterestId())
	if err != nil {
		return &pb.AddPointOfInterestResponse{
			Status: "Failed",
		}, err
	}

	group, err := s.roommateStore.GetGroup(ctx, req.GetGroupId())
	if err != nil {
		return &pb.AddPointOfInterestResponse{
			Status: "Failed",
		}, err
	}

	newGroupPOI := &models.RoommateGroupPointOfInterest{
		Address:           pointOfInterest.Address,
		RoommateId:        req.GetRoommateId(),
		PointOfInterestId: pointOfInterest.ID,
		PlaceId:           pointOfInterest.PlaceId,
	}

	var groupPOIs []*models.RoommateGroupPointOfInterest
	if group.PointsOfInterest == nil {
		groupPOIs = make([]*models.RoommateGroupPointOfInterest, 1)
	}
	updatedGroupPOIs := append(groupPOIs, newGroupPOI)
	group.PointsOfInterest = updatedGroupPOIs
	_, err = s.roommateStore.UpdateGroup(ctx, group)
	if err != nil {
		return &pb.AddPointOfInterestResponse{
			Status: "Failed",
		}, err
	}

	return &pb.AddPointOfInterestResponse{
		Status: "Updated",
	}, nil
}

func (s *RoommateService) GetGroupPointsOfInterest(ctx context.Context, req *pb.GetGroupPOIsRequest) (*pb.GetGroupPOIsResponse, error) {
	// TODO: validation
	groupId := req.GetGroupId()
	group, err := s.roommateStore.GetGroup(ctx, groupId)
	if err != nil {
		// TODO: return safe error message
		return &pb.GetGroupPOIsResponse{}, err
	}
	var pointsOfInterest []*pb.PointOfInterest
	for _, poi := range group.PointsOfInterest {
		pointsOfInterest = append(pointsOfInterest, &pb.PointOfInterest{
			Id:      poi.ID,
			PlaceId: poi.PlaceId,
		})
	}
	return &pb.GetGroupPOIsResponse{
		PointsOfInterest: pointsOfInterest,
	}, nil
}
