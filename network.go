package byte_ns

import (
	"container/heap"
	"go.uber.org/atomic"
	"runtime"
)

// Network Indicates a simulated network, which contains some simulated nodes
type Network struct {
	nodes          []Node
	running        *atomic.Bool
	runningCount   *atomic.Int32
	loopLimit      int
	emptySpinLimit int
	splitThreshold int
}

func NewNetwork(nodes []Node, loopLimit, emptySpinLimit, splitThreshold int) *Network {
	return &Network{
		nodes:          nodes,
		running:        atomic.NewBool(false),
		runningCount:   atomic.NewInt32(0),
		loopLimit:      loopLimit,
		emptySpinLimit: emptySpinLimit,
		splitThreshold: splitThreshold,
	}
}

// fetch Fetch Packets from nodes in the network, and put them into given heap
func (n *Network) fetch(packetHeap heap.Interface) bool {
	flag := false
	for _, node := range n.nodes {
		buffer := node.Packets()
		if buffer == nil {
			continue
		}
		buffer.Reduce(func(packet *SimulatedPacket) {
			heap.Push(packetHeap, packet)
			flag = true
		})
	}
	return flag
}

// drain Drain the given heap if possible, and OnEmit the Packets available
func (n *Network) drain(packetHeap *PacketHeap) bool {
	flag := false
	t := Now()
	for !packetHeap.IsEmpty() {
		p := packetHeap.Peek()
		if p.EmitTime.After(t) {
			break
		}
		p.Where.Emit(p)
		heap.Pop(packetHeap)
		flag = true
	}
	return flag
}

func (n *Network) clear(packetHeap *PacketHeap) {
	for !packetHeap.IsEmpty() {
		n.drain(packetHeap)
	}
}

func (n *Network) split(packetHeap *PacketHeap) {
	count := n.runningCount.Inc()
	if int(count) <= n.loopLimit {
		length := packetHeap.Len()
		h := &PacketHeap{packetHeap.storage[length/2:]}
		packetHeap.storage = packetHeap.storage[:length/2]
		heap.Init(packetHeap)
		heap.Init(h)
		go n.mainLoop(h, count)
	} else {
		n.runningCount.Dec()
	}
}

// mainLoop Main polling loop of network
func (n *Network) mainLoop(packetHeap *PacketHeap, index int32) {
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
	go n.mainLoop(&PacketHeap{}, 1)
}

// Stop the network, release resources
func (n *Network) Stop() {
	n.running.Store(false)
}
