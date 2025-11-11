package service

import (
	"context"
	"io"
	"sync"
	"time"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"homesearch.axel.to/commute/internal/clients"
	"homesearch.axel.to/commute/internal/store"
	"homesearch.axel.to/commute/internal/util"
	pb "homesearch.axel.to/services/commute/api"
	"homesearch.axel.to/shared/logger"
)

const (
	matrixFieldMask = "destinationIndex,duration"
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
	computeRequest, err := util.BuildComputeMatrixRequest(rentalPlaceId, poiPlaceIds, request.GetDisableTraffic())
	if err != nil {
		logger.LogAttrs(ctx, logger.LevelError, "Failed to build compute matrix request", logger.Group("parameters", logger.String("rentalPlaceId", rentalPlaceId), logger.Any("pointsOfInterestPlaceIds", poiPlaceIds), logger.Bool("disableTraffic", request.GetDisableTraffic())), logger.Any("error", err))
		return &pb.CalculateResponse{}, status.Error(codes.Internal, "Failed to compute commute")
	}

	computeCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	computeCtx = metadata.AppendToOutgoingContext(computeCtx, "X-Goog-Fieldmask", matrixFieldMask)
	defer cancel()

	logger.LogAttrs(computeCtx, logger.LevelInfo, "send compute matrix request to google api", logger.Any("request", computeRequest))
	resultsStream, err := cs.gRoutesService.ComputeRouteMatrix(computeCtx, computeRequest)
	if err != nil {
		logger.LogAttrs(computeCtx, logger.LevelError, "Failed to get route service matrix stream", logger.Any("error", err))
		return &pb.CalculateResponse{}, status.Error(codes.Internal, "Failed to receive commute durations")
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
		go processor.ProcessRouteElement(routeElement)
	}
	group.Wait()

	if len(streamErrors) > 0 {
		logger.LogAttrs(ctx, logger.LevelError, "errors during commute stream", logger.Any("errors", streamErrors))
	}

	// TODO: process errors here and send safe non-specific messages to client
	results, errs := processor.GetResults()
	if len(errs) > 0 {
		logger.LogAttrs(ctx, logger.LevelError, "errors during commute processing", logger.Any("errors", errs))
	}

	// TODO: store commute result in commute store if there are no errors
	response := &pb.CalculateResponse{
		PointsOfInterest: results,
	}
	logger.LogAttrs(ctx, logger.LevelInfo, "response", logger.Any("response", response))
	return response, nil
}
