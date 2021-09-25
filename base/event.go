package base

import (
	"time"
)

// Event is an action to be acted at the exactly time point
type Event interface {
	// Time when to act
	Time() time.Time
	// Action what to do
	Action() Action
	// HookBefore hook the event with the given action handled before the actual action
	HookBefore(action Action)
	// HookAfter hook the event with the given action handled after the actual action
	HookAfter(action Action)
}

type event struct {
	time   time.Time
	action Action
}

// Action is what to do of an event, time of the event is passed in, return following events of this event
type Action func(time.Time) []Event

// RepeatAction is same to Action, but repeated after the return delay.
// if return delay is negative, no longer repeat
type RepeatAction func(time.Time) ([]Event, time.Duration)

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) Action() Action {
	return e.action
}

func (e *event) HookBefore(action Action) {
	actualAction := e.action
	e.action = func(t time.Time) (events []Event) {
		for _, event := range action(t) {
			events = append(events, event)
		}
		for _, event := range actualAction(t) {
			events = append(events, event)
		}
		return
	}
}

func (e *event) HookAfter(action Action) {
	actualAction := e.action
	e.action = func(t time.Time) (events []Event) {
		for _, event := range actualAction(t) {
			events = append(events, event)
		}
		for _, event := range action(t) {
			events = append(events, event)
		}
		return
	}
}

// Aggregate the events into a slice, utility function
func Aggregate(events ...Event) []Event {
	return events
}

// NewDelayedEvent create an event with delay
func NewDelayedEvent(action Action, delay time.Duration, now time.Time) Event {
	return NewFixedEvent(action, now.Add(delay))
}

// NewFixedEvent create an event at the time point
func NewFixedEvent(action Action, time time.Time) Event {
	return &event{time: time, action: action}
}

// NewPeriodicEvent create a periodic event, generate itself each time
func NewPeriodicEvent(action Action, period time.Duration, t time.Time) Event {
	return NewRepeatEvent(func(now time.Time) ([]Event, time.Duration) {
		return action(now), period
	}, t)
}

// NewRepeatEvent create an event repeated after the return delay, if such delay is not negative
func NewRepeatEvent(action RepeatAction, t time.Time) Event {
	var actualAction func(now time.Time) []Event
	actualAction = func(now time.Time) []Event {
		events, delay := action(now)
		if delay >= 0 {
			events = append(events, NewDelayedEvent(actualAction, delay, now))
		}
		return events
	}
	return NewFixedEvent(actualAction, t)
}
