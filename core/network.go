package core

import (
	"container/heap"
	"go.uber.org/atomic"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes   []Node
	running *atomic.Bool
}

func NewNetwork(nodes []Node) *Network {
	return &Network{nodes: nodes, running: atomic.NewBool(false)}
}

// fetch Fetch Packets from nodes in the network, and put them into given heap
func (n *Network) fetch(packetHeap heap.Interface) {
	for _, node := range n.nodes {
		buffer := node.Packets()
		if buffer == nil {
			continue
		}
		buffer.Reduce(func(packet *SimulatedPacket) {
			heap.Push(packetHeap, packet)
		})
	}
}

// drain Drain the given heap if possible, and Emit the Packets available
func (n *Network) drain(packetHeap *PacketHeap) {
	t := Now()
	for !packetHeap.IsEmpty() {
		p := packetHeap.Peek()
		if p.EmitTime.After(t) {
			break
		}
		p.Where.Emit(p)
		heap.Pop(packetHeap)
	}
}

// mainLoop Main polling loop of network
func (n *Network) mainLoop() {
	if !n.running.CAS(false, true) {
		return
	}
	println("network main loop start")
	packetHeap := &PacketHeap{}
	for n.running.Load() {
		n.fetch(packetHeap)
		n.drain(packetHeap)
	}
	for !packetHeap.IsEmpty() {
		n.drain(packetHeap)
	}
	println("network main loop end")
}

// Start the network to enable packet transmission
func (n *Network) Start() {
	go n.mainLoop()
}

// Stop the network, release resources
func (n *Network) Stop() {
	n.running.Store(false)
}
