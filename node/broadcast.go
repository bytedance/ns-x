package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// BroadcastNode is a node broadcast any packet pass by to all its next node
// Although some other nodes have the same behavior, broadcast node is designed as a transparent node to avoid side effects.
type BroadcastNode struct {
	*BasicNode
}

// NewBroadcastNode creates a new BroadcastNode with given options
func NewBroadcastNode(options ...Option) *BroadcastNode {
	n := &BroadcastNode{
		BasicNode: &BasicNode{},
	}
	apply(n, options...)
	return n
}

func (n *BroadcastNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	events := make([]base.Event, len(n.GetNext()))
	for index, node := range n.GetNext() {
		events[index] = base.NewFixedEvent(func(t time.Time) []base.Event {
			return n.actualTransfer(packet, n, node, t)
		}, now)
	}
	return events
}
