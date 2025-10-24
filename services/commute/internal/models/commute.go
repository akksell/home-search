package models

import (
	"time"

	base "homesearch.axel.to/base/types"
)

type GroupCommute struct {
	ID                  string                    `firestore:"-"`
	GroupId             string                    `firestore:"groupId"`
	RentalHomeId        string                    `firestore:"rentalHomeId"`
	RentalHomeAddress   *base.Address             `firestore:"rentalHomeAddress"`
	PointsOfInterestIds []string                  `firestore:"pointsOfInterestIds"`
	PointsOfInterest    []*CommutePointOfInterest `firestore:"pointsOfInterest"`
	LastComputedAt      time.Time                 `firestore:"lastComputedAt"`
	CommuteParameters   *CommuteComputeParameters `firestore:"parameters"`
}

type CommutePointOfInterest struct {
	PointOfInterestLabel string        `firestore:"pointOfInterestLabel,omitempty"`
	Address              *base.Address `firestore:"address"`
	DurationSeconds      int32         `firestore:"durationSeconds"`
	IsTrafficAware       bool          `firestore:"isTrafficAware"`
}

type CommuteComputeParameters struct {
	DateOfCommute    time.Time `firestore:"dateOfCommute"`
	TimeOfDay        time.Time `firestore:"timeOfDay"`
	UTCOffsetSeconds int32     `firestore:"utcOffsetSeconds"`
}
