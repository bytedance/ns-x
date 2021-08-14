package node

import (
	"go.uber.org/atomic"
	"math"
	"ns-x/base"
	"time"
)

// RestrictNode simulate a node with limited ability
// Once Events through a RestrictNode reaches the limit(in bps or pps), the later Events will be put in a events
// Once the events overflow, later Events will be discarded
// The events limit will not be accurate, usually a little lower than specified, since it takes time(usually less than microseconds) to send Events
type RestrictNode struct {
	*BasicNode
	ppsLimit, bpsLimit                float64
	bufferSizeLimit, bufferCountLimit int64
	bufferSize, bufferCount           *atomic.Int64
	busyTime                          *atomic.Int64 // unix nano sec
}

// NewRestrictNode create a new restrict with the given parameter
// next, recordSize, callback the same as BasicNode
// ppsLimit, bpsLimit: the limit of Events per second/bytes per second
// bufferSizeLimit, bufferCountLimit: the limit of waiting Events, in bytes/Events
func NewRestrictNode(name string, onEmitCallback base.TransferCallback, ppsLimit, bpsLimit float64, bufferSizeLimit, bufferCountLimit int64) *RestrictNode {
	return &RestrictNode{
		BasicNode:        NewBasicNode(name, onEmitCallback),
		ppsLimit:         ppsLimit,
		bpsLimit:         bpsLimit,
		bufferSizeLimit:  bufferSizeLimit,
		bufferCountLimit: bufferCountLimit,
		bufferSize:       atomic.NewInt64(0),
		bufferCount:      atomic.NewInt64(0),
		busyTime:         atomic.NewInt64(0),
	}
}

func (n *RestrictNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	// TODO: an accurate way to determine buffer overflow
	if n.bufferSize.Load() >= n.bufferSizeLimit || n.bufferCount.Load() >= n.bufferCountLimit {
		return nil
	}
	flag := true
	for flag {
		nanoseconds := n.busyTime.Load()
		t := time.Unix(0, nanoseconds)
		flag = false
		if t.Before(now) {
			flag = !n.busyTime.CAS(nanoseconds, now.UnixNano())
		}
	}
	step := math.Max(1.0/n.ppsLimit, float64(packet.Size())/n.bpsLimit)
	delta := (time.Duration(step*1000*1000) * time.Microsecond).Nanoseconds()
	t := time.Unix(0, n.busyTime.Add(delta)-delta)
	n.bufferSize.Add(int64(packet.Size()))
	n.bufferCount.Inc()
	return base.Aggregate(
		base.NewFixedEvent(func() []base.Event {
			n.bufferSize.Sub(int64(packet.Size()))
			n.bufferCount.Dec()
			return n.ActualEmit(packet, n.next[0], t)
		}, t),
	)
}

func (n *RestrictNode) Check() {
	if n.next == nil || len(n.next) != 1 {
		panic("restrict node can only has single connection")
	}
}
