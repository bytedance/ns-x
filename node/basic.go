package node

import (
	"ns-x/base"
	"time"
)

// BasicNode is skeleton implementation of Node
type BasicNode struct {
	name           string
	next           []base.Node
	events         *base.EventBuffer
	onEmitCallback base.OnEmitCallback
}

// NewBasicNode creates a new BasicNode
func NewBasicNode(name string, onEmitCallback base.OnEmitCallback) *BasicNode {
	return &BasicNode{
		name:           name,
		next:           []base.Node{},
		events:         base.NewEventBuffer(),
		onEmitCallback: onEmitCallback,
	}
}

func (n *BasicNode) Name() string {
	return n.name
}

func (n *BasicNode) Events() *base.EventBuffer {
	return n.events
}

func (n *BasicNode) ActualEmit(packet base.Packet, target base.Node, now time.Time) {
	target.Emit(packet, now)
	if n.onEmitCallback != nil {
		n.onEmitCallback(packet, target, now)
	}
}

func (n *BasicNode) Emit(base.Packet, time.Time) {
	panic("not implemented")
}

func (n *BasicNode) GetNext() []base.Node {
	return n.next
}

func (n *BasicNode) SetNext(nodes ...base.Node) {
	n.next = nodes
}

func (n *BasicNode) Check() {
}
