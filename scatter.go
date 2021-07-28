package byte_ns

type PathSelector func(packet *SimulatedPacket, record *PacketQueue, nodes []Node) Node

type Scatter struct {
	BasicNode
	selector PathSelector
}

func NewScatter(selector PathSelector) *Scatter {
	return &Scatter{
		BasicNode: BasicNode{},
		selector:  selector,
	}
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
	s.OnSend(p)
	path := s.selector(p, s.record, s.next)
	if path != nil {
		path.Send(packet)
	}
	s.OnEmit(p)
}
