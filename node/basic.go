package node

import (
	"ns-x/base"
)

// BasicNode is skeleton implementation of Node
type BasicNode struct {
	name           string
	next           []base.Node
	buffer         *base.PacketBuffer
	record         *base.PacketQueue
	onEmitCallback base.OnEmitCallback
}

// NewBasicNode creates a new BasicNode
func NewBasicNode(name string, recordSize int, onEmitCallback base.OnEmitCallback) *BasicNode {
	return &BasicNode{
		name:           name,
		next:           []base.Node{},
		buffer:         base.NewPacketBuffer(),
		record:         base.NewPacketQueue(recordSize),
		onEmitCallback: onEmitCallback,
	}
}

func (n *BasicNode) Name() string {
	return n.name
}

// OnSend do some common tasks, should be called by Send of implementations
func (n *BasicNode) OnSend(packet *base.SimulatedPacket) {
	if n.record != nil {
		n.record.Enqueue(packet)
	}
}

func (n *BasicNode) Packets() *base.PacketBuffer {
	return n.buffer
}

// OnEmit do some common tasks, should be called by Emit of implementations
func (n *BasicNode) OnEmit(packet *base.SimulatedPacket) {
	if n.onEmitCallback != nil {
		n.onEmitCallback(packet)
	}
}

func (n *BasicNode) Emit(packet *base.SimulatedPacket) {
	n.OnEmit(packet)
	if packet.Loss {
		return
	}
	for _, node := range n.next {
		node.Send(packet.Actual)
	}
}

func (n *BasicNode) Send(packet []byte) {
	panic("not implemented")
}

func (n *BasicNode) GetNext() []base.Node {
	return n.next
}

func (n *BasicNode) SetNext(nodes ...base.Node) {
	if nodes == nil {
		panic("cannot set next nodes to nil")
	}
	n.next = nodes
}
