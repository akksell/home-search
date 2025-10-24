package models

import (
	"time"

	base "homesearch.axel.to/base/types"
)

type WrappedAddress struct {
	ID                    string               `firestore:"-"`
	CurrentPlaceId        string               `firestore:"currentPlaceId"`
	Address               *base.Address        `firestore:"addressComponents"`
	AddressHash           string               `firestore:"addressHash"`
	InitialCreateDate     time.Time            `firestore:"initialCreateDate"`
	LastPlaceIdUpdateDate time.Time            `firestore:"lastPlaceIdUpdateDate"`
	LastRefreshDate       time.Time            `firestore:"lastRefreshDate"`
	LastThreePlaceIds     []*DeprecatedPlaceId `firestore:"lastThreePlaceIds"`
}

type DeprecatedPlaceId struct {
	PlaceId                 string    `firestore:"placeId"`
	DeprecatedAt            time.Time `firestore:"deprecatedAt"`
	OriginallyValidFromDate time.Time `firestore:"originallyValidFromDate"`
}
