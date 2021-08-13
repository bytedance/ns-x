package node

import (
	"ns-x/base"
	"time"
)

// RouteSelector is rule of how to select the path of packet
type RouteSelector func(packet base.Packet, nodes []base.Node) base.Node

// ScatterNode transfer packet pass by to one of its next nodes according to a given rule
type ScatterNode struct {
	*BasicNode
	selector RouteSelector
}

func NewScatterNode(name string, selector RouteSelector, callback base.OnEmitCallback) *ScatterNode {
	return &ScatterNode{
		BasicNode: NewBasicNode(name, callback),
		selector:  selector,
	}
}

func (s *ScatterNode) Emit(packet base.Packet, now time.Time) {
	path := s.selector(packet, s.next)
	if path != nil {
		s.Events().Insert(base.NewFixedEvent(func() {
			s.ActualEmit(packet, path, now)
		}, now))
	}
}
