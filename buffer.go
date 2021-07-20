package networksimulator

import (
	"sync/atomic"
	"unsafe"
)

type PacketBuffer struct {
	node *unsafe.Pointer
}

type node struct {
	next *node
	data *SimulatedPacket
}

func (b *PacketBuffer) Insert(packet *SimulatedPacket) {
	n := &node{next: (*node)(*b.node), data: packet}
	for !atomic.CompareAndSwapPointer(b.node, unsafe.Pointer(n.next), unsafe.Pointer(n)) {
		n.next = (*node)(*b.node)
	}
}

func (b *PacketBuffer) Reduce(action func(packet *SimulatedPacket)) {
	n := *b.node
	for !atomic.CompareAndSwapPointer(b.node, n, nil) {
		n = *b.node
	}
	node := (*node)(n)
	for node != nil {
		action(node.data)
		node = node.next
	}
}
