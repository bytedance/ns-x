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

func NewEndpointNode(name string, callback base.TransferCallback) *EndpointNode {
	return &EndpointNode{
		BasicNode: NewBasicNode(name, callback),
	}
}

func (n *EndpointNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	if n.callback != nil {
		n.callback(packet, now)
	}
	return nil
}

func (n EndpointNode) SendAt(packet base.Packet, t time.Time) base.Event {
	return base.NewFixedEvent(func(t time.Time) []base.Event {
		return n.ActualTransfer(packet, n.next[0], t)
	}, t)
}

func (n *EndpointNode) Send(packet base.Packet) base.Event {
	return n.SendAt(packet, time.Now())
}

func (n *EndpointNode) Receive(callback func(packet base.Packet, now time.Time)) {
	n.callback = callback
}

func (n *EndpointNode) Check() {
	if len(n.next) > 1 {
		panic("endpoint node can has at most single connection")
	}
	n.BasicNode.Check()
}
