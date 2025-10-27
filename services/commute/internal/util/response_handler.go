package util

import (
	"fmt"
	"sync"

	routingpb "cloud.google.com/go/maps/routing/apiv2/routingpb"
	"google.golang.org/genproto/googleapis/rpc/code"

	base "homesearch.axel.to/base/types"
	pb "homesearch.axel.to/services/commute/api"
	roommatepb "homesearch.axel.to/services/roommate/api"
)

const (
	hoursInDay      = 24
	minutesInHour   = 60
	secondsInMinute = 60
	secondsInDay    = secondsInMinute * minutesInHour * hoursInDay
	secondsInHour   = secondsInMinute * minutesInHour
)

type calculateResponseResult struct {
	results          []*pb.PointOfInterestDuration
	processingErrors []error
}

type calculateResponseProcessor struct {
	lock             *sync.Mutex
	group            *sync.WaitGroup
	pointsOfInterest []*roommatepb.PointOfInterest
	processedResults *calculateResponseResult
}

func NewCalculateResponseProcessor(lock *sync.Mutex, group *sync.WaitGroup, groupPointsOfInterest []*roommatepb.PointOfInterest) *calculateResponseProcessor {
	return &calculateResponseProcessor{
		lock:             lock,
		group:            group,
		pointsOfInterest: groupPointsOfInterest,
		processedResults: &calculateResponseResult{
			results:          make([]*pb.PointOfInterestDuration, 0),
			processingErrors: make([]error, 0),
		},
	}
}

func (p *calculateResponseProcessor) ProcessRouteElement(routeElement *routingpb.RouteMatrixElement) {
	defer p.group.Done()
	if routeElement.GetStatus().GetCode() != int32(code.Code_OK) {
		processorError := fmt.Errorf("Failed to process point of interest %d: status %d with message %s", routeElement.GetDestinationIndex(), int32(routeElement.GetStatus().GetCode()), routeElement.GetStatus().GetMessage())
		p.lock.Lock()
		p.processedResults.processingErrors = append(p.processedResults.processingErrors, processorError)
		p.lock.Unlock()
		return
	}

	poiIndex := routeElement.GetDestinationIndex()
	pointOfInterest := p.pointsOfInterest[poiIndex]
	commuteDuration := routeElement.GetDuration()
	commuteSeconds := commuteDuration.GetSeconds()

	days, commuteSeconds := computeDurationComponent(commuteSeconds, secondsInDay)
	hours, commuteSeconds := computeDurationComponent(commuteSeconds, secondsInHour)
	minutes, commuteSeconds := computeDurationComponent(commuteSeconds, secondsInMinute)

	result := &pb.PointOfInterestDuration{
		PointOfInterest: pointOfInterest,
		Duration: &base.Duration{
			Days:    days,
			Hours:   hours,
			Minutes: minutes,
			Seconds: uint32(commuteSeconds),
		},
	}

	p.lock.Lock()
	p.processedResults.results = append(p.processedResults.results, result)
	p.lock.Unlock()
}

func (p *calculateResponseProcessor) GetResults() ([]*pb.PointOfInterestDuration, []error) {
	return p.processedResults.results, p.processedResults.processingErrors
}

func computeDurationComponent(commuteInSeconds, componentBoundary int64) (uint32, int64) {
	var updatedCommuteInSeconds int64 = commuteInSeconds
	var component uint32 = 0
	if commuteInSeconds/componentBoundary > 0 {
		component = uint32(commuteInSeconds / componentBoundary)
		updatedCommuteInSeconds = commuteInSeconds % componentBoundary
	}
	return component, updatedCommuteInSeconds
}
