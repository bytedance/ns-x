package byte_ns

import (
	"go.uber.org/atomic"
	"unsafe"
)

// PacketBuffer is ...
type PacketBuffer struct {
	node *atomic.UnsafePointer
}

// NewPacketBuffer creates a new packet buffer
func NewPacketBuffer() *PacketBuffer {
	return &PacketBuffer{
		node: atomic.NewUnsafePointer(nil),
	}
}

type node struct {
	next *node
	data *SimulatedPacket
}

// Insert ...
func (b *PacketBuffer) Insert(packet *SimulatedPacket) {
	n := &node{next: (*node)(b.node.Load()), data: packet}
	for !b.node.CAS(unsafe.Pointer(n.next), unsafe.Pointer(n)) {
		n.next = (*node)(b.node.Load())
	}
}

// Reduce ...
func (b *PacketBuffer) Reduce(action func(packet *SimulatedPacket)) {
	n := b.node.Swap(nil)
	node := (*node)(n)
	for node != nil {
		action(node.data)
		node = node.next
	}
}
