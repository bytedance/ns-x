package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// RouteSelector is rule of how to select the path of packet
type RouteSelector func(packet base.Packet, nodes []base.Node) base.Node

// ScatterNode transfer packet pass by to one of its next nodes according to a given rule
type ScatterNode struct {
	*BasicNode
	selector RouteSelector
}

// NewScatterNode create a ScatterNode with given options
func NewScatterNode(options ...Option) *ScatterNode {
	n := &ScatterNode{
		BasicNode: &BasicNode{},
	}
	apply(n, options...)
	if n.selector == nil {
		panic("a route selector must be specified for scatter nodes")
	}
	return n
}

func (n *ScatterNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	path := n.selector(packet, n.GetNext())
	if path != nil {
		return base.Aggregate(
			base.NewFixedEvent(func(t time.Time) []base.Event {
				return n.actualTransfer(packet, n, path, t)
			}, now),
		)
	}
	return nil
}

// WithRouteSelector create an option to set/overwrite route selector of nodes applied
// The nodes applied must be a ScatterNode
func WithRouteSelector(selector RouteSelector) Option {
	return func(node base.Node) {
		n, ok := node.(*ScatterNode)
		if !ok {
			panic("cannot set route selector")
		}
		n.selector = selector
	}
}
