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

func (b *Broadcast) Send(packet *Packet) {
	for _, n := range b.nodes {
		n.Send(packet)
	}
	t := Now()
	p := &SimulatedPacket{
		Actual:   packet,
		EmitTime: t,
		SentTime: t,
		Loss:     false,
		Where:    b,
	}
	b.BasicNode.OnSend(p)
}
