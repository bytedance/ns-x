package byte_ns

// Gather ...
type Gather struct {
	BasicNode
}

func NewGather(name string) *Gather {
	return &Gather{BasicNode{name: name}}
}

func (g *Gather) Send(packet *Packet) {
	t := Now()
	p := &SimulatedPacket{
		Actual:   packet,
		SentTime: t,
		EmitTime: t,
		Where:    g,
		Loss:     false,
	}
	g.OnSend(p)
	g.Emit(p)
}
