package node

import (
	"ns-x/base"
	"time"
)

// EndpointNode is a node to receive Events, Events reached an endpoint will no longer be transmitted
type EndpointNode struct {
	*BasicNode
	callback func(packet base.Packet, now time.Time)
}

func NewEndpointNode(name string, callback base.OnEmitCallback) *EndpointNode {
	return &EndpointNode{
		BasicNode: NewBasicNode(name, callback),
	}
}

func (n *EndpointNode) Emit(packet base.Packet, now time.Time) {
	if n.callback != nil {
		n.callback(packet, now)
	}
}

func (n *EndpointNode) Send(packet base.Packet) {
	now := time.Now()
	for _, node := range n.next {
		n.Events().Insert(base.NewFixedEvent(func() {
			n.ActualEmit(packet, node, now)
		}, now))
	}
}

func (n *EndpointNode) Receive(callback func(packet base.Packet, now time.Time)) {
	n.callback = callback
}
