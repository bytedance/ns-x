package node

import (
	"network-simulator/core"
)

type PathSelector func(packet *core.SimulatedPacket, record *core.PacketQueue, nodes []core.Node) core.Node

type Scatter struct {
	BasicNode
	selector PathSelector
	paths    []core.Node
}

func (s *Scatter) Send(packet *core.Packet) {
	t := core.Now()
	p := &core.SimulatedPacket{
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
