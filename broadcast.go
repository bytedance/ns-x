package networksimulator

type Broadcast struct {
	BasicNode
	nodes []Node
}

func NewBroadcast(nodes []Node, callback OnEmitCallback) *Broadcast {
	return &Broadcast{
		BasicNode: *NewBasicNode(nil, 0, callback),
		nodes:     nodes,
	}
}

func (b *Broadcast) Emit(packet *SimulatedPacket) {
	for _, n := range b.nodes {
		n.Send(packet.Actual)
	}
	b.BasicNode.Emit(packet)
}
