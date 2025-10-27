package models

import (
	"time"
)

type PointOfInterest struct {
	ID            string    `firestore:"-"`
	PlaceId       string    `firestore:"placeId"`
	RoommateId    string    `firestore:"roommateId"`
	Label         string    `firestore:"label,omitempty"`
	Tags          []string  `firestore:"tags,omitempty"`
	Address       *Address  `firestore:"address"`
	CreatedAt     time.Time `firestore:"createdAt,omitempty"`
	LastUpdatedAt time.Time `firestore:"lastUpdatedAt,omitempty"`
}

type Address struct {
	Street            string `firestore:"street"`
	City              string `firestore:"city"`
	PostalCode        string `firestore:"postalCode"`
	StateProvinceCode string `firestore:"stateProvinceCode,omitempty"`
	CountryCode       string `firestore:"countryCode"`
}
