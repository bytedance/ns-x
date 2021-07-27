package byte_ns

type PathSelector func(packet *SimulatedPacket, record *PacketQueue, nodes []Node) Node

type Scatter struct {
	BasicNode
	selector PathSelector
	paths    []Node
}

func (s *Scatter) Send(packet *Packet) {
	t := Now()
	p := &SimulatedPacket{
		Actual:   packet,
		EmitTime: t,
		SentTime: t,
		Loss:     false,
		Where:    s,
	}
	path := s.selector(p, s.record, s.paths)
	if path != nil {
		path.Send(packet)
	}
	s.BasicNode.OnSend(p)
}
