package is

import (
	"time"

	"github.com/patrickward/datacop"
)

// Before checks if a time is before another
func Before(t time.Time) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(time.Time)
		if !ok {
			return false
		}
		return v.Before(t)
	}
}

// After checks if a time is after another
func After(t time.Time) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(time.Time)
		if !ok {
			return false
		}
		return v.After(t)
	}
}

// BetweenTime checks if a value is between two other values
//
// Example usage:
// BetweenTime(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))(time.Now()) // returns true
// BetweenTime(time.Now().Add(-1*time.Hour), time.Now().Add(-30*time.Minute))(time.Now()) // returns false
func BetweenTime(start, end time.Time) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(time.Time)
		if !ok {
			return false
		}
		return v.After(start) && v.Before(end)
	}
}
