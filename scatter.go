package byte_ns

// PathSelector is rule of how to select the path of packet
type PathSelector func(packet *SimulatedPacket, record *PacketQueue, nodes []Node) Node

// Scatter transfer packet pass by to one of its next nodes according to a given rule
type Scatter struct {
	BasicNode
	selector PathSelector
}

func NewScatter(name string, selector PathSelector) *Scatter {
	return &Scatter{
		BasicNode: BasicNode{name: name},
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
