package byte_ns

// Broadcast is ...
type Broadcast struct {
	BasicNode
	nodes []Node
}

// NewBroadcast creates a new Broadcast with given Node(s)
func NewBroadcast(nodes []Node, callback OnEmitCallback) *Broadcast {
	return &Broadcast{
		BasicNode: *NewBasicNode(0, callback),
		nodes:     nodes,
	}
}

func (b *Broadcast) Emit(packet *SimulatedPacket) {
	for _, n := range b.nodes {
		n.Send(packet.Actual)
	}
	b.BasicNode.Emit(packet)
}
