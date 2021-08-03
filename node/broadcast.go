package node

import (
	"byte-ns/base"
)

// Broadcast is a node broadcast any packet pass by to all its next node
// Although some other nodes have the same behavior, broadcast node is designed as a transparent node to avoid side effects.
type Broadcast struct {
	BasicNode
}

// NewBroadcast creates a new Broadcast with given Node(s)
func NewBroadcast(name string, callback base.OnEmitCallback) *Broadcast {
	return &Broadcast{
		BasicNode: *NewBasicNode(name, 0, callback),
	}
}

func (b *Broadcast) Send(packet *base.Packet) {
	for _, n := range b.next {
		n.Send(packet)
	}
	t := base.Now()
	p := &base.SimulatedPacket{
		Actual:   packet,
		EmitTime: t,
		SentTime: t,
		Loss:     false,
		Where:    b,
	}
	b.BasicNode.OnSend(p)
}
