package models

import (
	"time"

	base "homesearch.axel.to/base/types"
)

type PointOfInterest struct {
	ID            string        `firestore:"-"`
	PlaceId       string        `firestore:"placeId"`
	RoommateId    string        `firestore:"roommateId"`
	Label         string        `firestore:"label,omitempty"`
	Tags          []string      `firestore:"tags,omitempty"`
	Address       *base.Address `firestore:"address"`
	CreatedAt     time.Time     `firestore:"createdAt,omitempty"`
	LastUpdatedAt time.Time     `firestore:"lastUpdatedAt,omitempty"`
}
