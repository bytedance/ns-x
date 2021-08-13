package base

import "time"

// Node Indicates a simulated node in the network
type Node interface {
	Name() string
	Events() *EventBuffer
	Emit(packet Packet, now time.Time)
	GetNext() []Node
	SetNext(nodes ...Node)
	Check()
}

// OnEmitCallback called when a packet is emitted
type OnEmitCallback func(packet Packet, target Node, now time.Time)
