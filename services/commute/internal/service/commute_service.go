package service

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"google.golang.org/grpc/metadata"

	"homesearch.axel.to/commute/internal/clients"
	"homesearch.axel.to/commute/internal/store"
	"homesearch.axel.to/commute/internal/util"
	pb "homesearch.axel.to/services/commute/api"
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
		return nil, fmt.Errorf("Failed to compute commute duration: missing group id")
	}

	homeAddress := request.GetHomeAddress()
	placeIdResponse, err := cs.addressWrapperService.GetPlaceId(ctx, homeAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to get address from address wrapper service: %w", err)
	}
	rentalPlaceId := placeIdResponse.GetPlaceId()
	// TODO: try and fetch a commute from the store using the placeId and hashed address
	//		before fetching from google api. If found, return early

	groupId := request.GetGroupId()
	wrappedGroupPOIs, err := cs.roommateService.GetGroupPointsOfInterest(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get group points of interest from roommate service: %w", err)
	}
	groupPOIs := wrappedGroupPOIs.GetPointsOfInterest()
	if len(groupPOIs) == 0 {
		return nil, fmt.Errorf("Failed to compute commute durations: points of interest for group id %s is empty", groupId)
	}

	// TODO: pull points of interest from firestore db, don't have time rn but
	// 		 will implement soon™. For now, these are all manual placeIds
	// tempPlaceIds := []string{"ChIJdbYmq-vTTYYRu_ZHX7-WoUU"}
	poiPlaceIds := make([]string, 0)
	for _, poi := range groupPOIs {
		poiPlaceIds = append(poiPlaceIds, poi.PlaceId)
	}
	computeRequest, err := util.BuildComputeMatrixRequest(rentalPlaceId, poiPlaceIds, request.GetDisableTraffic())
	if err != nil {
		return nil, err
	}
	fmt.Printf("built compute matrix request\n")

	computeCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	computeCtx = metadata.AppendToOutgoingContext(computeCtx, "X-Goog-Fieldmask", matrixFieldMask)
	defer cancel()

	resultsStream, err := cs.gRoutesService.ComputeRouteMatrix(computeCtx, computeRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to get route service matrix stream: %w", err)
	}

	var streamErrors []error
	lock := sync.Mutex{}
	group := sync.WaitGroup{}
	processor := util.NewCalculateResponseProcessor(&lock, &group, groupPOIs)
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

	// TODO: process errors here and send safe non-specific messages to client
	results, _ := processor.GetResults()

	// TODO: store commute result in commute store if there are no errors

	return &pb.CalculateResponse{
		PointsOfInterest: results,
	}, nil
}
