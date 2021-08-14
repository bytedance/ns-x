package base

import (
	"time"
)

type Event interface {
	Time() time.Time
	Action() func() []Event
}

type event struct {
	time   time.Time
	action func() []Event
}

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) Action() func() []Event {
	return e.action
}

func Aggregate(events ...Event) []Event {
	return events
}

func NewDelayedEvent(action func(time.Time) []Event, delay time.Duration, now time.Time) Event {
	t := now.Add(delay)
	return NewFixedEvent(func() []Event {
		return action(t)
	}, t)
}

func NewFixedEvent(action func() []Event, time time.Time) Event {
	return &event{time: time, action: action}
}
