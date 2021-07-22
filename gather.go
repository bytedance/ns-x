package networksimulator

type Gather struct {
	BasicNode
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
	g.buffer.Insert(p)
	g.BasicNode.Send(p)
}
