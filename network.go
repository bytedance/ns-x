package networksimulator

import (
	"container/heap"
	"go.uber.org/atomic"
	"time"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes   []Node
	running atomic.Bool
}

// fetch Fetch packets from nodes in the network, and put them into given heap
func (n *Network) fetch(packetHeap heap.Interface) {
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
	for p != nil && t.Before(p.EmitTime) {
		p.Where.emit(p)
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
