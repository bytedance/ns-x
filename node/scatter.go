package node

import (
	"github.com/bytedance/ns-x/base"
	"time"
)

// RouteSelector is rule of how to select the path of packet
type RouteSelector func(packet base.Packet, nodes []base.Node) base.Node

// ScatterNode transfer packet pass by to one of its next nodes according to a given rule
type ScatterNode struct {
	*BasicNode
	selector RouteSelector
}

func NewScatterNode(name string, selector RouteSelector, callback base.TransferCallback) *ScatterNode {
	return &ScatterNode{
		BasicNode: NewBasicNode(name, callback),
		selector:  selector,
	}
}

func (s *ScatterNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	path := s.selector(packet, s.next)
	if path != nil {
		return base.Aggregate(
			base.NewFixedEvent(func(t time.Time) []base.Event {
				return s.ActualTransfer(packet, path, t)
			}, now),
		)
	}
	return nil
}
