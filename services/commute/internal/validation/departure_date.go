package validation

import (
	"fmt"

	base "homesearch.axel.to/base/types"
)

const (
	maxNanoseconds = 0
	maxSeconds     = 59
	maxMinutes     = 59
	maxHours       = 23
)

func ValidateDepartureTime(timezoneTime *base.TimezoneAwareTime) error {
	time := timezoneTime.GetTime()
	if time == nil {
		return fmt.Errorf("Departure is missing a time")
	}
	if time.GetHour() > maxHours {
		return fmt.Errorf("Departure hours are greater than maximum allowed")
	}
	if time.GetMinute() > maxMinutes {
		return fmt.Errorf("Departure minutes are greater than maximum allowed")
	}
	if time.GetSecond() > maxSeconds {
		return fmt.Errorf("Departure seconds are greater than maxmimum allowed")
	}
	if time.GetNanosecond() > maxNanoseconds {
		return fmt.Errorf("Departure nanoseconds are greater than maximum allowed")
	}
	return nil
}
