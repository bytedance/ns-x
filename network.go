package ns_x

import (
	"container/heap"
	"go.uber.org/atomic"
	"ns-x/base"
	"runtime"
	"sync"
	"time"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes   []base.Node
	buffer  *base.EventBuffer
	running *atomic.Bool
	wg      *sync.WaitGroup
}

// NewNetwork creates a network with the given nodes, connections of nodes should be already established.
// loopLimit is the limit of parallelized main loops count
// a main loop will exit once spun emptySpinLimit rounds doing nothing
// a main loop will try to split into two loops once count of packets waiting to handle reach splitThreshold
func NewNetwork(nodes []base.Node) *Network {
	return &Network{
		nodes:   nodes,
		buffer:  base.NewEventBuffer(),
		running: atomic.NewBool(false),
		wg:      &sync.WaitGroup{},
	}
}

// fetch events from nodes in the network, and put them into given heap
func (n *Network) fetch(packetHeap heap.Interface) {
	n.buffer.Reduce(func(packet base.Event) {
		heap.Push(packetHeap, packet)
	})
}

// drain the given heap if possible, and process the Events available
func (n *Network) drain(packetHeap *base.EventHeap) {
	now := time.Now()
	for !packetHeap.IsEmpty() {
		p := packetHeap.Peek()
		t := p.Time()
		if t.After(now) {
			break
		}
		events := p.Action()(t)
		heap.Pop(packetHeap)
		for _, event := range events {
			heap.Push(packetHeap, event)
		}
	}
}

func (n *Network) clear(packetHeap *base.EventHeap) {
	for !packetHeap.IsEmpty() {
		n.drain(packetHeap)
	}
}

// eventLoop Main polling loop of network
func (n *Network) eventLoop(packetHeap *base.EventHeap) {
	println("network main loop start")
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	n.wg.Add(1)
	defer n.wg.Done()
	for n.running.Load() {
		n.fetch(packetHeap)
		n.drain(packetHeap)
	}
	n.clear(packetHeap)
	println("network main loop end")
}

// Start the network to enable packet transmission
func (n *Network) Start(events ...base.Event) {
	if n.running.Load() {
		return
	}
	n.running.Store(true)
	for _, node := range n.nodes {
		node.Check()
	}
	h := &base.EventHeap{Storage: events}
	heap.Init(h)
	go n.eventLoop(h)
}

// Stop the network, release resources
func (n *Network) Stop() {
	n.running.Store(false)
	n.wg.Wait()
}

func (n *Network) Event(events ...base.Event) {
	n.buffer.Insert(events...)
}
