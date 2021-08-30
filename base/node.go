package base

import "time"

// Node indicates a simulated node in the network
type Node interface {
	// Check whether the node can work correctly, usually called by network just before the simulation
	Check()
	// GetNext nodes of the node
	GetNext() []Node
	// SetNext nodes of the node, should not be used during simulation
	SetNext(nodes ...Node)
	// Transfer the given packet to somewhere at sometime the node decides
	// time calculation should use the given time as current time point
	// return the following events caused by the packet transfer
	Transfer(packet Packet, now time.Time) []Event
	// GetTransferCallback get the TransferCallback of the node
	GetTransferCallback() TransferCallback
	// SetTransferCallback set the TransferCallback of the node, should not be used during simulation
	SetTransferCallback(callback TransferCallback)
}

// TransferCallback called when a packet is transferred
type TransferCallback func(packet Packet, source, target Node, now time.Time)
