package networksimulator

import (
	"container/heap"
	"go.uber.org/atomic"
	"net"
	"time"
)

// Packet Indicates an Actual packet, with its data and address
type Packet struct {
	data    []byte
	address net.Addr
}

// SimulatedPacket Indicates a simulated packet, with its Actual packet and some simulated environment
type SimulatedPacket struct {
	Actual   *Packet   // the Actual packet
	EmitTime time.Time // when this packet is emitted (Where emit a packet means the packet leaves the Where, send to the next Where)
	SentTime time.Time // when this packet is sent (Where send a packet means the packet enters the Where, waiting to emit)
	Loss     bool      // whether this packet is lost
	Where    Node      // where is the packet
}

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

// Node Indicates a simulated node in the network
type Node interface {
	Send(packet *Packet) // Send a packet to the node
	packets() *PacketBuffer
	emit(packet *SimulatedPacket)
}

// OnEmitCallback called when a packet is emitted
type OnEmitCallback func(packet *SimulatedPacket)

// BasicNode implement Node basically
type BasicNode struct {
	next           Node
	buffer         *PacketBuffer
	record         *PacketQueue
	onEmitCallback OnEmitCallback
}

func NewBasicNode(next Node, recordSize int, onEmitCallback OnEmitCallback) *BasicNode {
	return &BasicNode{
		next:           next,
		buffer:         &PacketBuffer{},
		record:         NewPacketQueue(recordSize),
		onEmitCallback: onEmitCallback,
	}
}

func (n *BasicNode) Send(packet *SimulatedPacket) {
	n.record.Enqueue(packet)
}

func (n *BasicNode) packets() *PacketBuffer {
	return n.buffer
}

func (n *BasicNode) emit(packet *SimulatedPacket) {
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
