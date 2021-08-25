package tick

import "time"

type Clock func() time.Time

func NewRealClock() Clock {
	return time.Now
}

func NewStepClock(t time.Time, step time.Duration) Clock {
	return func() time.Time {
		result := t
		t = t.Add(step)
		return result
	}
}
