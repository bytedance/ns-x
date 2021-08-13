package base

import (
	"time"
)

type Event interface {
	Time() time.Time
	Action() func()
}

type event struct {
	time   time.Time
	action func()
}

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) Action() func() {
	return e.action
}

func NewDelayedEvent(action func(time.Time), delay time.Duration, now time.Time) Event {
	t := now.Add(delay)
	return NewFixedEvent(func() {
		action(t)
	}, t)
}

func NewFixedEvent(action func(), time time.Time) Event {
	return &event{time: time, action: action}
}
