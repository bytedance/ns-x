package base

import "time"

// Node Indicates a simulated node in the network
type Node interface {
	Name() string
	Transfer(packet Packet, now time.Time) []Event
	GetNext() []Node
	SetNext(nodes ...Node)
	Check()
}

// TransferCallback called when a packet is emitted
type TransferCallback func(packet Packet, target Node, now time.Time)
