package util

import (
	routingpb "cloud.google.com/go/maps/routing/apiv2/routingpb"
)

func BuildComputeMatrixRequest(rentalPlaceId string, pointsOfInterestPlaceIds []string, disableTraffic bool) (*routingpb.ComputeRouteMatrixRequest, error) {
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

	var routingPreference routingpb.RoutingPreference
	if disableTraffic {
		routingPreference = routingpb.RoutingPreference_TRAFFIC_UNAWARE
	} else {
		routingPreference = routingpb.RoutingPreference_TRAFFIC_AWARE_OPTIMAL
	}

	// ignore error since we know the values are ok
	winterHolidaySeason, _ := NewDateRange("2025-12-09", "2026-01-15")
	summerHolidaySeason, _ := NewDateRange("2026-05-26", "2026-09-01")
	// I feel like this sucks since I don't want to export the type
	// but oh well I'll improve this later
	exlcludedMonths := []dateRange{winterHolidaySeason, summerHolidaySeason}
	generator, err := NewDateGenerator("07:00:00", exlcludedMonths)
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
		TravelMode:        routingpb.RouteTravelMode_DRIVE,
		RoutingPreference: routingPreference,
		DepartureTime:     departureTime,
	}, nil
}
