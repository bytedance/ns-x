package node

import (
	"byte-ns/base"
	"byte-ns/time"
)

// PathSelector is rule of how to select the path of packet
type PathSelector func(packet *base.SimulatedPacket, record *base.PacketQueue, nodes []base.Node) base.Node

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

func (s *Scatter) Send(packet *base.Packet) {
	t := time.Now()
	p := &base.SimulatedPacket{
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
