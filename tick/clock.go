package tick

import "time"

// Clock always return current time in its view
type Clock func() time.Time

// NewRealClock create a clock align to system clock
func NewRealClock() Clock {
	return time.Now
}

// NewStepClock create a clock, return the calculated time and step forward a given duration once called
func NewStepClock(t time.Time, step time.Duration) Clock {
	return func() time.Time {
		result := t
		t = t.Add(step)
		return result
	}
}
