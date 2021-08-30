package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"math"
	"time"
)

// RestrictNode simulate a node with limited ability
// Once packets through a RestrictNode reaches the limit(in bps or pps), the later packets will be put in a buffer
// Once the buffer overflow, later packets will be discarded
type RestrictNode struct {
	*BasicNode
	ppsLimit, bpsLimit                float64
	bufferSizeLimit, bufferCountLimit int64
	bufferSize, bufferCount           int64
	busyTime                          time.Time
}

// NewRestrictNode create a new RestrictNode with the given options
func NewRestrictNode(options ...Option) *RestrictNode {
	n := &RestrictNode{
		BasicNode:        &BasicNode{},
		ppsLimit:         -1,
		bpsLimit:         -1,
		bufferSizeLimit:  -1,
		bufferCountLimit: -1,
	}
	apply(n, options...)
	if n.ppsLimit <= 0 && n.bpsLimit <= 0 {
		panic("a restrict node must be limited in pps/bps")
	}
	return n
}

func (n *RestrictNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	delay := false
	t := now
	if n.busyTime.After(now) {
		t = n.busyTime
		delay = true
	}
	if delay {
		if n.bufferSizeLimit >= 0 && n.bufferSize+int64(packet.Size()) >= n.bufferSizeLimit {
			return nil
		}
		if n.bufferCountLimit >= 0 && n.bufferCount+1 >= n.bufferCountLimit {
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
			return n.actualTransfer(packet, n, n.GetNext()[0], t)
		}, t),
	)
}

func (n *RestrictNode) Check() {
	if len(n.GetNext()) != 1 {
		panic("restrict node can only has single connection")
	}
}

// WithPPSLimit create an option set/overwrite pps limit and buffer count limit to nodes applied
// once flow of the node calculated in packets/second reach pps limit, further packets will be put into the buffer
// once total count of packets in the buffer reach the buffer count limit, further packets will be ignored
// node applied must be a RestrictNode
// set limit to -1 means unlimited
func WithPPSLimit(ppsLimit float64, bufferCountLimit int64) Option {
	return func(node base.Node) {
		n, ok := node.(*RestrictNode)
		if !ok {
			panic("cannot set pps limit")
		}
		n.ppsLimit = ppsLimit
		n.bufferCountLimit = bufferCountLimit
	}
}

// WithBPSLimit create an option set/overwrite bps limit and buffer size limit to nodes applied
// once flow of the node calculated in bytes/second reach pps limit, further packets will be put into the buffer
// once total size of packets in the buffer reach the buffer size limit, further packets will be ignored
// node applied must be a RestrictNode
// set limit to -1 means unlimited
func WithBPSLimit(bpsLimit float64, bufferSizeLimit int64) Option {
	return func(node base.Node) {
		n, ok := node.(*RestrictNode)
		if !ok {
			panic("cannot set pps limit")
		}
		n.bpsLimit = bpsLimit
		n.bufferSizeLimit = bufferSizeLimit
	}
}
