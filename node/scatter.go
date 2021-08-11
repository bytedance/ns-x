package node

import (
	"ns-x/base"
	"ns-x/time"
)

// RouteSelector is rule of how to select the path of packet
type RouteSelector func(packet *base.SimulatedPacket, record *base.PacketQueue, nodes []base.Node) base.Node

// ScatterNode transfer packet pass by to one of its next nodes according to a given rule
type ScatterNode struct {
	BasicNode
	selector RouteSelector
}

func NewScatterNode(name string, selector RouteSelector) *ScatterNode {
	return &ScatterNode{
		BasicNode: BasicNode{name: name},
		selector:  selector,
	}
}

func (s *ScatterNode) Send(packet []byte) {
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
