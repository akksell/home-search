package util

import (
	"strconv"

	routingpb "cloud.google.com/go/maps/routing/apiv2/routingpb"

	base "homesearch.axel.to/base/types"
	"homesearch.axel.to/commute/internal/validation"
	pb "homesearch.axel.to/services/commute/api"
)

type RequestType uint8

const (
	REQUEST_TYPE_HOME_DEPARTURE RequestType = iota
	REQUEST_TYPE_POI_DEPARTURE
)

func BuildComputeMatrixRequest(request *pb.CalculateRequest, rentalPlaceId string, pointsOfInterestPlaceIds []string, requestType RequestType) (*routingpb.ComputeRouteMatrixRequest, error) {
	var routingPreference routingpb.RoutingPreference
	if request.GetDisableTraffic() {
		routingPreference = routingpb.RoutingPreference_TRAFFIC_UNAWARE
	} else {
		routingPreference = routingpb.RoutingPreference_TRAFFIC_AWARE_OPTIMAL
	}
	travelMode := routingpb.RouteTravelMode_DRIVE

	// ignore error since we know the values are ok
	winterHolidaySeason, _ := NewDateRange("2025-12-09", "2026-01-15")
	summerHolidaySeason, _ := NewDateRange("2026-05-26", "2026-09-01")
	// I feel like this sucks since I don't want to export the type
	// but oh well I'll improve this later
	excludedMonths := []dateRange{winterHolidaySeason, summerHolidaySeason}
	if requestType == REQUEST_TYPE_HOME_DEPARTURE {
		return buildHomeDepartureRequest(
			request,
			rentalPlaceId,
			pointsOfInterestPlaceIds,
			excludedMonths,
			routingPreference,
			travelMode,
		)
	}
	return buildPOIDepartureRequest(
		request,
		rentalPlaceId,
		pointsOfInterestPlaceIds,
		excludedMonths,
		routingPreference,
		travelMode,
	)
}

func buildHomeDepartureRequest(
	request *pb.CalculateRequest,
	rentalPlaceId string,
	pointsOfInterestPlaceIds []string,
	excludedMonths []dateRange,
	routingPreference routingpb.RoutingPreference,
	travelMode routingpb.RouteTravelMode,
) (*routingpb.ComputeRouteMatrixRequest, error) {
	homeDepartureTime := request.GetHomeDepartureTime()
	homeDepartureTimeString := "07:00:00"
	if homeDepartureTime != nil {
		err := validation.ValidateDepartureTime(homeDepartureTime)
		if err != nil {
			return nil, err
		}
		homeDepartureTimeString = timezoneAwareTimeToString(homeDepartureTime)
	}

	rentalOrigin := &routingpb.RouteMatrixOrigin{
		Waypoint: &routingpb.Waypoint{
			LocationType: &routingpb.Waypoint_PlaceId{
				PlaceId: rentalPlaceId,
			},
		},
	}
	origins := []*routingpb.RouteMatrixOrigin{rentalOrigin}

	var destinations []*routingpb.RouteMatrixDestination
	for _, id := range pointsOfInterestPlaceIds {
		destination := &routingpb.RouteMatrixDestination{
			Waypoint: &routingpb.Waypoint{
				LocationType: &routingpb.Waypoint_PlaceId{
					PlaceId: id,
				},
			},
		}
		destinations = append(destinations, destination)
	}

	generator, err := NewDateGenerator(homeDepartureTimeString, excludedMonths)
	if err != nil {
		return nil, err
	}
	departureTime, err := generator.GenerateTimestamp()
	if err != nil {
		return nil, err
	}

	return &routingpb.ComputeRouteMatrixRequest{
		Origins:           origins,
		Destinations:      destinations,
		TravelMode:        travelMode,
		RoutingPreference: routingPreference,
		DepartureTime:     departureTime,
	}, nil
}

func buildPOIDepartureRequest(
	request *pb.CalculateRequest,
	rentalPlaceId string,
	pointsOfInterestPlaceIds []string,
	excludedMonths []dateRange,
	routingPreference routingpb.RoutingPreference,
	travelMode routingpb.RouteTravelMode,
) (*routingpb.ComputeRouteMatrixRequest, error) {
	poiDepartureTime := request.GetPointOfInterestDepartureTime()
	poiDepartureTimeString := "16:30:00"
	if poiDepartureTime != nil {
		err := validation.ValidateDepartureTime(poiDepartureTime)
		if err != nil {
			return nil, err
		}
		poiDepartureTimeString = timezoneAwareTimeToString(poiDepartureTime)
	}

	rentalDestination := &routingpb.RouteMatrixDestination{
		Waypoint: &routingpb.Waypoint{
			LocationType: &routingpb.Waypoint_PlaceId{
				PlaceId: rentalPlaceId,
			},
		},
	}
	destinations := []*routingpb.RouteMatrixDestination{rentalDestination}

	var origins []*routingpb.RouteMatrixOrigin
	for _, id := range pointsOfInterestPlaceIds {
		origin := &routingpb.RouteMatrixOrigin{
			Waypoint: &routingpb.Waypoint{
				LocationType: &routingpb.Waypoint_PlaceId{
					PlaceId: id,
				},
			},
		}
		origins = append(origins, origin)
	}

	generator, err := NewDateGenerator(poiDepartureTimeString, excludedMonths)
	if err != nil {
		return nil, err
	}
	departureTime, err := generator.GenerateTimestamp()
	if err != nil {
		return nil, err
	}

	return &routingpb.ComputeRouteMatrixRequest{
		Origins:           origins,
		Destinations:      destinations,
		TravelMode:        travelMode,
		RoutingPreference: routingPreference,
		DepartureTime:     departureTime,
	}, nil
}

func timezoneAwareTimeToString(timezoneTime *base.TimezoneAwareTime) string {
	base := 10
	hours := strconv.FormatUint(uint64(timezoneTime.GetTime().GetHour()), base)
	minutes := strconv.FormatUint(uint64(timezoneTime.GetTime().GetMinute()), base)
	seconds := strconv.FormatUint(uint64(timezoneTime.GetTime().GetSecond()), base)
	return padLeft(hours, 2, '0') + ":" + padLeft(minutes, 2, '0') + ":" + padLeft(seconds, 2, '0')
}

func padLeft(s string, length int, char rune) string {
	paddedString := s
	for len(paddedString) < length {
		paddedString = string(char) + paddedString
	}
	return paddedString
}
