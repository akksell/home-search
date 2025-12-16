package service

import (
	"context"
	"io"
	"sync"
	"time"

	routing "cloud.google.com/go/maps/routing/apiv2"
	routingpb "cloud.google.com/go/maps/routing/apiv2/routingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"homesearch.axel.to/commute/internal/clients"
	"homesearch.axel.to/commute/internal/store"
	"homesearch.axel.to/commute/internal/util"
	pb "homesearch.axel.to/services/commute/api"
	roommatepb "homesearch.axel.to/services/roommate/api"
	"homesearch.axel.to/shared/logger"
)

const (
	matrixFieldMask = "destinationIndex,duration,originIndex"
)

type commuteService struct {
	pb.UnimplementedCommuteServiceServer
	addressWrapperService *clients.AddressWrapperServiceClient
	roommateService       *clients.RoommateServiceClient
	gRoutesService        *routing.RoutesClient
	store                 *store.CommuteStore
}

func NewCommuteService(addressWrapperSvc *clients.AddressWrapperServiceClient, roommateSvc *clients.RoommateServiceClient, gRoutesSvc *routing.RoutesClient, commuteStore *store.CommuteStore) *commuteService {
	return &commuteService{
		addressWrapperService: addressWrapperSvc,
		roommateService:       roommateSvc,
		gRoutesService:        gRoutesSvc,
		store:                 commuteStore,
	}
}

func (cs *commuteService) Calculate(ctx context.Context, request *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	// TODO: validate request
	if request.GetGroupId() == "" {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to compute commute duration: missing group id")
		return &pb.CalculateResponse{}, status.Error(codes.InvalidArgument, "missing group id")
	}

	homeAddress := request.GetHomeAddress()
	placeIdResponse, err := cs.addressWrapperService.GetPlaceId(ctx, homeAddress)
	if err != nil {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to get address from address wrapper service", logger.Any("error", err))
		return &pb.CalculateResponse{}, status.Error(codes.Internal, "Failed to resolve address")
	}
	rentalPlaceId := placeIdResponse.GetPlaceId()
	// TODO: try and fetch a commute from the store using the placeId and hashed address
	//		before fetching from google api. If found, return early

	groupId := request.GetGroupId()
	wrappedGroupPOIs, err := cs.roommateService.GetGroupPointsOfInterest(ctx, groupId)
	if err != nil {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to get group points of interest from roommate service", logger.Any("error", err))
		return &pb.CalculateResponse{}, status.Error(codes.Internal, "Failed to get group points of interest")
	}
	groupPOIs := wrappedGroupPOIs.GetPointsOfInterest()
	if len(groupPOIs) == 0 {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to compute commute durations - points of interest are empty", logger.Group("group", logger.String("id", groupId), logger.Any("pointsOfInterest", groupPOIs)))
		return &pb.CalculateResponse{}, status.Error(codes.InvalidArgument, "Group has 0 points of interest")
	}

	poiPlaceIds := make([]string, 0)
	for _, poi := range groupPOIs {
		poiPlaceIds = append(poiPlaceIds, poi.PlaceId)
	}
	homeDepartureComputeRequest, err := util.BuildComputeMatrixRequest(request, rentalPlaceId, poiPlaceIds, util.REQUEST_TYPE_HOME_DEPARTURE)
	if err != nil {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to build compute matrix request", logger.Group("parameters", logger.String("rentalPlaceId", rentalPlaceId), logger.Any("pointsOfInterestPlaceIds", poiPlaceIds), logger.Bool("disableTraffic", request.GetDisableTraffic())), logger.Any("error", err))
		return &pb.CalculateResponse{}, status.Error(codes.Internal, "Failed to compute commute")
	}

	poiDepartureComputeRequest, err := util.BuildComputeMatrixRequest(request, rentalPlaceId, poiPlaceIds, util.REQUEST_TYPE_POI_DEPARTURE)
	if err != nil {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to build compute matrix request", logger.Group("parameters", logger.String("rentalPlaceId", rentalPlaceId), logger.Any("pointsOfInterestPlaceIds", poiPlaceIds), logger.Bool("disableTraffic", request.GetDisableTraffic())), logger.Any("error", err))
		return &pb.CalculateResponse{}, status.Error(codes.Internal, "Failed to compute commute")
	}

	homeDepartureResults, homeDepartureErrs := cs.sendCommuteRequest(ctx, homeDepartureComputeRequest, groupPOIs, util.REQUEST_TYPE_HOME_DEPARTURE)
	poiDepartureResults, poiDepartureErrs := cs.sendCommuteRequest(ctx, poiDepartureComputeRequest, groupPOIs, util.REQUEST_TYPE_POI_DEPARTURE)
	var processingErrors []error
	processingErrors = append(processingErrors, homeDepartureErrs...)
	processingErrors = append(processingErrors, poiDepartureErrs...)

	// TODO: process errors here and send safe non-specific messages to client
	if len(processingErrors) > 0 {
		logger.LogAttrs(ctx, logger.LevelError, "errors during commute processing", logger.Any("errors", processingErrors))
	}

	// TODO: store commute result in commute store if there are no errors
	response := &pb.CalculateResponse{
		DepartFrom: &pb.CalculateResponse_DepartReturn{
			PointsOfInterest: homeDepartureResults,
		},
		ReturnTo: &pb.CalculateResponse_DepartReturn{
			PointsOfInterest: poiDepartureResults,
		},
	}
	logger.LogAttrs(ctx, logger.LevelInfo, "response", logger.Any("response", response))
	return response, nil
}

func (cs *commuteService) sendCommuteRequest(ctx context.Context, computeRequest *routingpb.ComputeRouteMatrixRequest, groupPOIs []*roommatepb.PointOfInterest, requestType util.RequestType) ([]*pb.PointOfInterestDuration, []error) {
	computeCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	computeCtx = metadata.AppendToOutgoingContext(computeCtx, "X-Goog-Fieldmask", matrixFieldMask)
	defer cancel()

	logger.LogAttrs(computeCtx, logger.LevelInfo, "send compute matrix request to google api", logger.Any("request", computeRequest))
	resultsStream, err := cs.gRoutesService.ComputeRouteMatrix(computeCtx, computeRequest)
	if err != nil {
		logger.LogAttrs(computeCtx, logger.LevelError, "Failed to get route service matrix stream", logger.Any("error", err))
		return nil, []error{err}
	}

	var streamErrors []error
	lock := sync.Mutex{}
	group := sync.WaitGroup{}
	processor := util.NewCalculateResponseProcessor(&lock, &group, groupPOIs)
	logger.LogAttrs(ctx, logger.LevelInfo, "processing commute duration results")
	for {
		routeElement, err := resultsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			streamErrors = append(streamErrors, err)
			continue
		}
		group.Add(1)
		logger.LogAttrs(ctx, logger.LevelInfo, "processing routeElement", logger.Any("routeElement", routeElement))
		go processor.ProcessRouteElement(routeElement, requestType)
	}
	group.Wait()
	return processor.GetResults()
}
