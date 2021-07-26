package node

import (
	"network-simulator/core"
)

// BasicNode implement Node basically
type BasicNode struct {
	next           core.Node
	buffer         *core.PacketBuffer
	record         *core.PacketQueue
	onEmitCallback core.OnEmitCallback
}

func NewBasicNode(next core.Node, recordSize int, onEmitCallback core.OnEmitCallback) *BasicNode {
	return &BasicNode{
		next:           next,
		buffer:         core.NewPackerBuffer(),
		record:         core.NewPacketQueue(recordSize),
		onEmitCallback: onEmitCallback,
	}
}

func (n *BasicNode) OnSend(packet *core.SimulatedPacket) {
	n.record.Enqueue(packet)
}

func (n *BasicNode) Packets() *core.PacketBuffer {
	return n.buffer
}

func (n *BasicNode) Emit(packet *core.SimulatedPacket) {
	if n.onEmitCallback != nil {
		n.onEmitCallback(packet)
	}
	if packet.Loss {
		return
	}
	if n.next != nil {
		n.next.Send(packet.Actual)
	}
}

func (n *BasicNode) Send(packet *core.Packet) {
	panic("not implemented")
}
