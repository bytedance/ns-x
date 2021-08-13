package base

import (
	"go.uber.org/atomic"
	"unsafe"
)

// EventBuffer is a thread-safe, lock-free buffer used to store simulated packets, implemented like a single link list
type EventBuffer struct {
	node *atomic.UnsafePointer
}

// NewEventBuffer creates a new packet buffer
func NewEventBuffer() *EventBuffer {
	return &EventBuffer{
		node: atomic.NewUnsafePointer(nil),
	}
}

type node struct {
	next  *node
	event Event
}

// Insert a simulated packet to the buffer, thread-safe
func (b *EventBuffer) Insert(events ...Event) {
	for _, event := range events {
		n := &node{next: (*node)(b.node.Load()), event: event}
		for !b.node.CAS(unsafe.Pointer(n.next), unsafe.Pointer(n)) {
			n.next = (*node)(b.node.Load())
		}
	}
}

// Reduce means clear the buffer and do an action on the event cleared, thread-safe
func (b *EventBuffer) Reduce(action func(event Event)) {
	n := b.node.Swap(nil)
	node := (*node)(n)
	for node != nil {
		action(node.event)
		node = node.next
	}
}
