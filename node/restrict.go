package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"math"
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
	bufferSize, bufferCount           int64
	busyTime                          time.Time
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
	}
}

func (n *RestrictNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	delay := false
	t := now
	if n.busyTime.After(now) {
		t = n.busyTime
		delay = true
	}
	if delay {
		if n.bufferSize+int64(packet.Size()) >= n.bufferSizeLimit || n.bufferCount+1 >= n.bufferCountLimit {
			return nil
		}
	}
	step := math.Max(1.0/n.ppsLimit, float64(packet.Size())/n.bpsLimit)
	delta := time.Duration(step*1000*1000) * time.Microsecond
	n.busyTime = t.Add(delta)
	if delay {
		n.bufferSize += int64(packet.Size())
		n.bufferCount++
	}
	return base.Aggregate(
		base.NewFixedEvent(func(t time.Time) []base.Event {
			if delay {
				n.bufferSize -= int64(packet.Size())
				n.bufferCount--
			}
			return n.ActualTransfer(packet, n.next[0], t)
		}, t),
	)
}

func (n *RestrictNode) Check() {
	if len(n.next) != 1 {
		panic("restrict node can only has single connection")
	}
}
