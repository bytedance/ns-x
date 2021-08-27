package ns_x

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/tick"
	"go.uber.org/atomic"
	"runtime"
	"sync"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes   []base.Node
	clock   tick.Clock
	buffer  *base.EventBuffer
	running *atomic.Bool
	wg      *sync.WaitGroup
}

// NewNetwork creates a network with the given nodes, connections of nodes should be already established.
func NewNetwork(nodes []base.Node, clock tick.Clock) *Network {
	return &Network{
		nodes:   nodes,
		clock:   clock,
		buffer:  base.NewEventBuffer(),
		running: atomic.NewBool(false),
		wg:      &sync.WaitGroup{},
	}
}

// fetch events from nodes in the network, and put them into given heap
func (n *Network) fetch(eventQueue *base.EventQueue) {
	n.buffer.Reduce(func(event base.Event) {
		eventQueue.Enqueue(event)
	})
}

// drain the given heap if possible, and process the events available
func (n *Network) drain(eventQueue *base.EventQueue) {
	now := n.clock()
	for !eventQueue.IsEmpty() {
		p := eventQueue.Peek()
		t := p.Time()
		if t.After(now) {
			break
		}
		events := p.Action()(t)
		eventQueue.Dequeue()
		for _, event := range events {
			eventQueue.Enqueue(event)
		}
	}
}

// block until clear the given heap
func (n *Network) clear(eventQueue *base.EventQueue) {
	for !eventQueue.IsEmpty() {
		n.drain(eventQueue)
	}
}

// eventLoop Main polling loop of network
func (n *Network) eventLoop(eventQueue *base.EventQueue) {
	println("network main loop start")
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	n.wg.Add(1)
	defer n.wg.Done()
	for n.running.Load() {
		n.fetch(eventQueue)
		n.drain(eventQueue)
	}
	n.clear(eventQueue)
	println("network main loop end at", n.clock().String())
}

// Start the network to enable event process
func (n *Network) Start(config Config) {
	if n.running.Load() {
		return
	}
	n.running.Store(true)
	for _, node := range n.nodes {
		node.Check()
	}
	eventQueue := base.NewEventQueue(config.BucketSize, config.MaxBuckets)
	for _, event := range config.InitialEvents {
		eventQueue.Enqueue(event)
	}
	go n.eventLoop(eventQueue)
}

// Stop the network, release resources
func (n *Network) Stop() {
	n.running.Store(false)
	n.wg.Wait()
}

// Nodes return all nodes managed by the network
func (n *Network) Nodes() []base.Node {
	return n.nodes
}
