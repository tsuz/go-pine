package pine

import "time"

// TimeValue is a model for storing time and its associated value
type TimeValue struct {
	Time  time.Time
	Value float64
}

// NewTimeValue creates a value for time
func NewTimeValue(t time.Time, v float64) *TimeValue {
	return &TimeValue{
		Time:  t,
		Value: v,
	}
}
