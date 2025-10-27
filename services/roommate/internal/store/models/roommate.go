package models

import "time"

type Roommate struct {
	ID             string    `firestore:"-"`
	FirstName      string    `firestore:"firstName"`
	MiddleName     string    `firestore:"middleName,omitempty"`
	LastName       string    `firestore:"lastName"`
	MemberGroupIds []string  `firestore:"memberGroupIds"`
	CreatedAt      time.Time `firestore:"createdAt"`
	UpdatedAt      time.Time `firestore:"updatedAt"`
}
