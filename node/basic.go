package node

import (
	"github.com/bytedance/ns-x/base"
	"time"
)

// BasicNode is skeleton implementation of Node
type BasicNode struct {
	name     string
	next     []base.Node
	callback base.TransferCallback
}

// NewBasicNode creates a new BasicNode
func NewBasicNode(name string, callback base.TransferCallback) *BasicNode {
	return &BasicNode{
		name:     name,
		next:     []base.Node{},
		callback: callback,
	}
}

func (n *BasicNode) Name() string {
	return n.name
}

func (n *BasicNode) ActualTransfer(packet base.Packet, target base.Node, now time.Time) []base.Event {
	if n.callback != nil {
		n.callback(packet, target, now)
	}
	return target.Transfer(packet, now)
}

func (n *BasicNode) Transfer(base.Packet, time.Time) []base.Event {
	panic("not implemented")
	return nil
}

func (n *BasicNode) GetNext() []base.Node {
	return n.next
}

func (n *BasicNode) SetNext(nodes ...base.Node) {
	n.next = nodes
}

func (n *BasicNode) Check() {
}
