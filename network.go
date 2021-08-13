package ns_x

import (
	"container/heap"
	"go.uber.org/atomic"
	"ns-x/base"
	"runtime"
	"time"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes          []base.Node
	running        *atomic.Bool
	runningCount   *atomic.Int32
	loopLimit      int
	emptySpinLimit int
	splitThreshold int
}

// NewNetwork creates a network with the given nodes, connections of nodes should be already established.
// loopLimit is the limit of parallelized main loops count
// a main loop will exit once spun emptySpinLimit rounds doing nothing
// a main loop will try to split into two loops once count of packets waiting to handle reach splitThreshold
func NewNetwork(nodes []base.Node, loopLimit, emptySpinLimit, splitThreshold int) *Network {
	return &Network{
		nodes:          nodes,
		running:        atomic.NewBool(false),
		runningCount:   atomic.NewInt32(0),
		loopLimit:      loopLimit,
		emptySpinLimit: emptySpinLimit,
		splitThreshold: splitThreshold,
	}
}

// fetch Fetch Events from nodes in the network, and put them into given heap
func (n *Network) fetch(packetHeap heap.Interface) bool {
	flag := false
	for _, node := range n.nodes {
		buffer := node.Events()
		if buffer == nil {
			continue
		}
		buffer.Reduce(func(packet base.Event) {
			heap.Push(packetHeap, packet)
			flag = true
		})
	}
	return flag
}

// drain Drain the given heap if possible, and OnEmit the Events available
func (n *Network) drain(packetHeap *base.EventHeap) bool {
	flag := false
	t := time.Now()
	for !packetHeap.IsEmpty() {
		p := packetHeap.Peek()
		if p.Time().After(t) {
			break
		}
		p.Action()()
		heap.Pop(packetHeap)
		flag = true
	}
	return flag
}

func (n *Network) clear(packetHeap *base.EventHeap) {
	for !packetHeap.IsEmpty() {
		n.drain(packetHeap)
	}
}

func (n *Network) split(packetHeap *base.EventHeap) {
	count := n.runningCount.Inc()
	if int(count) <= n.loopLimit {
		length := packetHeap.Len()
		h := &base.EventHeap{Storage: packetHeap.Storage[length/2:]}
		packetHeap.Storage = packetHeap.Storage[:length/2]
		heap.Init(packetHeap)
		heap.Init(h)
		go n.eventLoop(h, count)
	} else {
		n.runningCount.Dec()
	}
}

// eventLoop Main polling loop of network
func (n *Network) eventLoop(packetHeap *base.EventHeap, index int32) {
	println("network main loop start #", index)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	emptySpinCount := 0
	for n.running.Load() {
		emptySpinCount++
		if n.fetch(packetHeap) {
			emptySpinCount = 0
		}
		if packetHeap.Len() > n.splitThreshold {
			n.split(packetHeap)
		}
		if n.drain(packetHeap) {
			emptySpinCount = 0
		}
		if emptySpinCount >= n.emptySpinLimit {
			count := n.runningCount.Dec()
			if count > 0 {
				n.clear(packetHeap)
				println("network main loop end #", index, "after spun", emptySpinCount, "rounds")
				return
			}
			n.runningCount.Inc()
		}
	}
	n.clear(packetHeap)
	n.runningCount.Dec()
	println("network main loop end #", index)
}

// Start the network to enable packet transmission
func (n *Network) Start() {
	if n.runningCount.Load() > 0 || n.running.Load() {
		return
	}
	n.running.Store(true)
	n.runningCount.Inc()
	for _, node := range n.nodes {
		node.Check()
	}
	go n.eventLoop(&base.EventHeap{}, 1)
}

// Stop the network, release resources
func (n *Network) Stop() {
	n.running.Store(false)
}
