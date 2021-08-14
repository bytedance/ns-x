package base

import (
	"go.uber.org/atomic"
	"time"
)

type Event interface {
	Time() time.Time
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

func Aggregate(events ...Event) []Event {
	return events
}

func NewDelayedEvent(action Action, delay time.Duration, now time.Time) Event {
	return NewFixedEvent(action, now.Add(delay))
}

func NewFixedEvent(action Action, time time.Time) Event {
	return &event{time: time, action: action}
}

type Cancel func()

func NewPeriodicEvent(action Action, period time.Duration, t time.Time) (Event, Cancel) {
	flag := atomic.NewBool(true)
	var actualAction func(now time.Time) []Event
	actualAction = func(now time.Time) []Event {
		events := action(now)
		if flag.Load() {
			events = append(events, NewFixedEvent(actualAction, now.Add(period)))
		}
		return events
	}
	return NewFixedEvent(actualAction, t), func() {
		flag.Store(false)
	}
}
