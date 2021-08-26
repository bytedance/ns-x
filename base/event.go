package base

import (
	"go.uber.org/atomic"
	"time"
)

// Event is an action to be acted at the exactly time point
type Event interface {
	// Time when to act
	Time() time.Time
	// Action what to do
	Action() Action
}

type event struct {
	time   time.Time
	action Action
}

type Action func(time.Time) []Event

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) Action() Action {
	return e.action
}

// Aggregate the events into a slice, utility function
func Aggregate(events ...Event) []Event {
	return events
}

// NewDelayedEvent create a event with delay
func NewDelayedEvent(action Action, delay time.Duration, now time.Time) Event {
	return NewFixedEvent(action, now.Add(delay))
}

// NewFixedEvent create a event at the time point
func NewFixedEvent(action Action, time time.Time) Event {
	return &event{time: time, action: action}
}

type Cancel func()

// NewPeriodicEvent create a periodic event, with a function to cancel it
func NewPeriodicEvent(action Action, period time.Duration, t time.Time) (Event, Cancel) {
	flag := atomic.NewBool(true)
	var actualAction func(now time.Time) []Event
	actualAction = func(now time.Time) []Event {
		if flag.Load() {
			events := action(now)
			events = append(events, NewFixedEvent(actualAction, now.Add(period)))
			return events
		}
		return nil
	}
	return NewFixedEvent(actualAction, t), func() {
		flag.Store(false)
	}
}
