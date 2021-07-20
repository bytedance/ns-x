package networksimulator

import (
	"container/heap"
	"go.uber.org/atomic"
	"net"
	"time"
)

type Packet struct {
	data    []byte
	address net.Addr
}

type SimulatedPacket struct {
	actual   *Packet
	emitTime time.Time
	sentTime time.Time
	loss     bool
	node     Node
}

type Network struct {
	nodes   []Node
	running atomic.Bool
}

// fetch Fetch packets from nodes in the network, and put them into given heap
func (n *Network) fetch(packetHeap *PacketHeap) {
	for _, node := range n.nodes {
		node.packets().Reduce(func(packet *SimulatedPacket) {
			heap.Push(packetHeap, packet)
		})
	}
}

// drain Drain the given heap if possible, and emit the packets available
func (n *Network) drain(packetHeap *PacketHeap) {
	t := time.Now()
	p := packetHeap.Peek()
	for p != nil && t.Before(p.emitTime) {
		p.node.emit(p)
		heap.Pop(packetHeap)
	}
}

// mainLoop Main polling loop of network
func (n *Network) mainLoop() {
	n.running.Store(true)
	packetHeap := &PacketHeap{}
	for n.running.Load() {
		n.fetch(packetHeap)
		n.drain(packetHeap)
	}
	for !packetHeap.IsEmpty() {
		n.drain(packetHeap)
	}
}

// Start the network to enable packet transmission
func (n *Network) Start() {
	go n.mainLoop()
}

// Stop the network, release resources
func (n *Network) Stop() {
	n.running.Store(false)
}

type Node interface {
	Send(packet *Packet)
	packets() *PacketBuffer
	emit(packet *SimulatedPacket)
}

type OnEmitCallback func(packet *SimulatedPacket)

type BasicNode struct {
	next           Node
	buffer         *PacketBuffer
	record         *PacketQueue
	onEmitCallback OnEmitCallback
}

func (n *BasicNode) Send(packet *SimulatedPacket) {
	n.record.Enqueue(packet)
}

func (n *BasicNode) packets() *PacketBuffer {
	return n.buffer
}

func (n *BasicNode) emit(packet *SimulatedPacket) {
	n.onEmitCallback(packet)
	if packet.loss {
		return
	}
	n.next.Send(packet.actual)
}
