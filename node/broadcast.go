package node

import "network-simulator/core"

type Broadcast struct {
	BasicNode
	nodes []core.Node
}

func NewBroadcast(nodes []core.Node, callback core.OnEmitCallback) *Broadcast {
	return &Broadcast{
		BasicNode: *NewBasicNode(nil, 0, callback),
		nodes:     nodes,
	}
}

func (b *Broadcast) Emit(packet *core.SimulatedPacket) {
	for _, n := range b.nodes {
		n.Send(packet.Actual)
	}
	b.BasicNode.Emit(packet)
}
