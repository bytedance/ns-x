package node

import (
	"ns-x/base"
	"ns-x/time"
)

// BroadcastNode is a node broadcast any packet pass by to all its next node
// Although some other nodes have the same behavior, broadcast node is designed as a transparent node to avoid side effects.
type BroadcastNode struct {
	BasicNode
}

// NewBroadcastNode creates a new BroadcastNode with given Node(s)
func NewBroadcastNode(name string, callback base.OnEmitCallback) *BroadcastNode {
	return &BroadcastNode{
		BasicNode: *NewBasicNode(name, 0, callback),
	}
}

func (b *BroadcastNode) Send(packet []byte) {
	for _, n := range b.next {
		n.Send(packet)
	}
	t := time.Now()
	p := &base.SimulatedPacket{
		Actual:   packet,
		EmitTime: t,
		SentTime: t,
		Loss:     false,
		Where:    b,
	}
	b.BasicNode.OnSend(p)
}
