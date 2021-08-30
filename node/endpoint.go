package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// EndpointNode is a node to send and receive packets
type EndpointNode struct {
	*BasicNode
	callback React
}

type React func(packet base.Packet, now time.Time) []base.Event

// NewEndpointNode create an EndpointNode with given options
func NewEndpointNode(options ...Option) *EndpointNode {
	n := &EndpointNode{
		BasicNode: &BasicNode{},
	}
	apply(n, options...)
	return n
}

func (n *EndpointNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	if n.callback != nil {
		return n.callback(packet, now)
	}
	return nil
}

func (n *EndpointNode) Send(packet base.Packet, t time.Time) base.Event {
	return base.NewFixedEvent(func(t time.Time) []base.Event {
		return n.actualTransfer(packet, n, n.GetNext()[0], t)
	}, t)
}

func (n *EndpointNode) Receive(callback React) {
	n.callback = callback
}

func (n *EndpointNode) Check() {
	if len(n.GetNext()) > 1 {
		panic("endpoint node can has at most single connection")
	}
	n.BasicNode.Check()
}
