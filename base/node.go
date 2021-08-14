package base

import "time"

// Node indicates a simulated node in the network
type Node interface {
	Name() string
	Check()
	GetNext() []Node
	SetNext(nodes ...Node)
	Transfer(packet Packet, now time.Time) []Event
}

// TransferCallback called when a packet is transferred
type TransferCallback func(packet Packet, target Node, now time.Time)
