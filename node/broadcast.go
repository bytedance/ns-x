package node

import (
	"ns-x/base"
	"time"
)

// BroadcastNode is a node broadcast any packet pass by to all its next node
// Although some other nodes have the same behavior, broadcast node is designed as a transparent node to avoid side effects.
type BroadcastNode struct {
	*BasicNode
}

// NewBroadcastNode creates a new BroadcastNode with given Node(s)
func NewBroadcastNode(name string, callback base.TransferCallback) *BroadcastNode {
	return &BroadcastNode{
		BasicNode: NewBasicNode(name, callback),
	}
}

func (n *BroadcastNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	events := make([]base.Event, len(n.next))
	for index, node := range n.next {
		events[index] = base.NewFixedEvent(func() []base.Event {
			return n.ActualEmit(packet, node, now)
		}, now)
	}
	return events
}
