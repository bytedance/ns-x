package node

import (
	"byte-ns/base"
	time2 "byte-ns/time"
	"go.uber.org/atomic"
	"math"
	"time"
)

// RestrictNode simulate a node with limited ability
// Once Packets through a RestrictNode reaches the limit(in bps or pps), the later Packets will be put in a buffer
// Once the buffer overflow, later Packets will be discarded
// The buffer limit will not be accurate, usually a little lower than specified, since it takes time(usually less than microseconds) to send Packets
type RestrictNode struct {
	BasicNode
	ppsLimit, bpsLimit                float64
	bufferSizeLimit, bufferCountLimit int64
	bufferSize, bufferCount           *atomic.Int64
	emitTime                          time.Time
}

// NewRestrictNode create a new restrict with the given parameter
// next, recordSize, onEmitCallback the same as BasicNode
// ppsLimit, bpsLimit: the limit of Packets per second/bytes per second
// bufferSizeLimit, bufferCountLimit: the limit of waiting Packets, in bytes/Packets
func NewRestrictNode(name string, recordSize int, onEmitCallback base.OnEmitCallback,
	ppsLimit, bpsLimit float64,
	bufferSizeLimit, bufferCountLimit int64) *RestrictNode {
	return &RestrictNode{
		BasicNode:        *NewBasicNode(name, recordSize, onEmitCallback),
		ppsLimit:         ppsLimit,
		bpsLimit:         bpsLimit,
		bufferSizeLimit:  bufferSizeLimit,
		bufferCountLimit: bufferCountLimit,
		bufferSize:       atomic.NewInt64(0),
		bufferCount:      atomic.NewInt64(0),
		emitTime:         time2.Now(),
	}
}

func (r *RestrictNode) Emit(packet *base.SimulatedPacket) {
	r.BasicNode.Emit(packet)
	r.bufferSize.Sub(int64(len(packet.Actual)))
	r.bufferCount.Dec()
}

func (r *RestrictNode) Send(packet []byte) {
	if r.bufferSize.Load() >= r.bufferSizeLimit || r.bufferCount.Load() >= r.bufferCountLimit {
		return
	}
	sentTime := time2.Now()
	if r.emitTime.Before(sentTime) {
		r.emitTime = sentTime
	}
	emitTime := r.emitTime
	p := &base.SimulatedPacket{Actual: packet, EmitTime: emitTime, SentTime: sentTime, Loss: false, Where: r}
	step := math.Max(1.0/r.ppsLimit, float64(len(packet))/r.bpsLimit)
	r.emitTime = emitTime.Add(time.Duration(step*1000*1000) * time.Microsecond)
	r.bufferSize.Add(int64(len(packet)))
	r.bufferCount.Inc()
	r.Packets().Insert(p)
	r.BasicNode.OnSend(p)
}
