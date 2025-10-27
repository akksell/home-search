package util

import (
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MONTHS_TO_ADD = 2
)

type dateRange struct {
	from time.Time
	to   time.Time
}

// Expects dates formatted as DateOnly time formatI
// see: https://pkg.go.dev/time@go1.25.3#pkg-constants
func NewDateRange(from, to string) (dateRange, error) {
	fromTime, err := time.Parse(time.DateOnly, from)
	if err != nil {
		return dateRange{}, err
	}
	toTime, err := time.Parse(time.DateOnly, to)
	if err != nil {
		return dateRange{}, err
	}
	if fromTime.After(toTime) {
		return dateRange{}, errors.New("Invalid Range: from is after to")
	}
	return dateRange{
		from: fromTime,
		to:   toTime,
	}, nil
}

// Only care about the months in the range, year is irrelevant
// it will be in the near future (< 1 year)
type dateGenerator struct {
	timeOfDay     time.Time
	excludeRanges []dateRange // this should be order and not overlapping
}

// Expects time formatted as TimeOnly time format
// see: https://pkg.go.dev/time@go1.25.3#pkg-constants
func NewDateGenerator(timeOfDay string, excludeRanges []dateRange) (dateGenerator, error) {
	timeOfDayAsTime, err := time.Parse(time.TimeOnly, timeOfDay)
	if err != nil {
		return dateGenerator{}, err
	}
	return dateGenerator{
		timeOfDay:     timeOfDayAsTime,
		excludeRanges: excludeRanges,
	}, nil
}

// Assumes the ranges passed are sequential, only considering month
// and day
func (dg dateGenerator) GenerateTimestamp() (*timestamppb.Timestamp, error) {
	// NOTE: Make sure to build with -tags timetzdata for this to work
	// TODO: should make this dynamic and relative to the users
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return nil, err
	}
	serverNow := time.Now().In(location)

	futureDate := serverNow.AddDate(0, MONTHS_TO_ADD, 0)
	monthsSeen := make([]time.Month, 0)
	for isWithinExcludedMonth(futureDate, dg.excludeRanges) {
		if len(monthsSeen) > 11 {
			break
		}
		monthsSeen = append(monthsSeen, futureDate.Month())
		futureDate = futureDate.AddDate(0, 1, 0)
	}
	// statistically busiest day to commute
	for futureDate.Weekday() != time.Thursday {
		futureDate = futureDate.AddDate(0, 0, 1)
	}

	// use 7AM as earliest (realisitic time we'd want to commute)
	// TODO: use a dynamic time value based on user preferences
	generatedTimestamp := time.Date(futureDate.Year(), futureDate.Month(), futureDate.Day(), 7, 0, 0, 0, location).UTC()
	return timestamppb.New(generatedTimestamp), nil
}

func isWithinExcludedMonth(timeValue time.Time, excludedRanges []dateRange) bool {
	for _, excludedRange := range excludedRanges {
		if excludedRange.from.Month() < timeValue.Month() &&
			excludedRange.to.Month() > timeValue.Month() {
			return true
		}
	}
	return false
}
