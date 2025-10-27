package models

import (
	"time"
)

type RoommateGroup struct {
	ID                       string                          `firestore:"-"`
	Name                     string                          `firestore:"name"`
	Description              string                          `firestore:"description,omitempty"`
	OwnerId                  string                          `firestore:"ownerId"`
	Owner                    *Roommate                       `firestore:"owner"`
	MemberIds                string                          `firestore:"memberIds,omitempty"`
	Members                  []*Roommate                     `firestore:"members,omitempty"`
	PointsOfInterest         []*RoommateGroupPointOfInterest `firestore:"pointsOfInterest,omitempty"`
	SelectedRentalPropertyId string                          `firestore:"selectedRentalId,omitempty"`
	JoinURL                  string                          `firestore:"joinURL"`
	CreatedAt                time.Time                       `firestore:"createdAt"`
	UpdatedAt                time.Time                       `firestore:"updatedAt"`
}

type RoommateGroupPointOfInterest struct {
	ID                string   `firestore:"-"`
	Address           *Address `firestore:"address"`
	PlaceId           string   `firestore:"placeId"`
	PointOfInterestId string   `firestore:"pointOfInterestId"`
	RoommateId        string   `firestore:"roommateId"`
}
