package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// BasicNode is skeleton implementation of Node
type BasicNode struct {
	next     []base.Node
	callback base.TransferCallback
}

func (n *BasicNode) actualTransfer(packet base.Packet, source, target base.Node, now time.Time) []base.Event {
	if n.callback != nil {
		n.callback(packet, source, target, now)
	}
	return target.Transfer(packet, now)
}

func (n *BasicNode) GetTransferCallback() base.TransferCallback {
	return n.callback
}

func (n *BasicNode) SetTransferCallback(callback base.TransferCallback) {
	n.callback = callback
}

func (n *BasicNode) GetNext() []base.Node {
	return n.next
}

func (n *BasicNode) SetNext(nodes ...base.Node) {
	n.next = nodes
}

func (n *BasicNode) Check() {
}

// WithTransferCallback create an option to set/overwrite the given transfer callback to nodes applied
// node applied must be a BasicNode
func WithTransferCallback(callback base.TransferCallback) Option {
	return func(node base.Node) {
		node.SetTransferCallback(callback)
	}
}
