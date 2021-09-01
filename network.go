package ns_x

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/tick"
	"runtime"
	"sync"
	"time"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes  []base.Node
	buffer *base.EventBuffer
	wg     *sync.WaitGroup
}

// NewNetwork creates a network with the given nodes, connections of nodes should be already established.
func NewNetwork(nodes []base.Node) *Network {
	return &Network{
		nodes:  nodes,
		buffer: base.NewEventBuffer(),
		wg:     &sync.WaitGroup{},
	}
}

// eventLoop Main polling loop of network
func (n *Network) eventLoop(eventQueue *base.EventQueue, clock tick.Clock, lifetime time.Duration) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer n.wg.Done()
	now := clock()
	deadline := now.Add(lifetime)
	println("network main loop start at", now.String())
	for !now.After(deadline) && !eventQueue.IsEmpty() {
		p := eventQueue.Peek()
		t := p.Time()
		if t.After(now) {
			now = clock()
			continue
		}
		events := p.Action()(t)
		eventQueue.Dequeue()
		for _, event := range events {
			eventQueue.Enqueue(event)
		}
	}
	println("network main loop end at", now.String())
}

// Run with the given config, users should Wait before another simulation or exit
// some Config can be used on the simulation, default valued will be used if not specified
// simulation will finish once no events remain or reach lifetime
func (n *Network) Run(events []base.Event, clock tick.Clock, lifetime time.Duration, configs ...Config) {
	n.wg.Add(1)
	for _, node := range n.nodes {
		node.Check()
	}
	config := &config{
		bucketSize: DefaultBucketSize,
		maxBuckets: DefaultMaxBuckets,
	}
	config.apply(configs...)
	eventQueue := base.NewEventQueue(config.bucketSize, config.maxBuckets)
	for _, event := range events {
		eventQueue.Enqueue(event)
	}
	go n.eventLoop(eventQueue, clock, lifetime)
}

// Wait until simulation finish
func (n *Network) Wait() {
	n.wg.Wait()
}

// Nodes return all nodes managed by the network
func (n *Network) Nodes() []base.Node {
	return n.nodes
}
