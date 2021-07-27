package networksimulator

import (
	"go.uber.org/atomic"
	"math"
	"time"
)

// Restrict simulate a node with limited ability
// Once Packets through a Restrict reaches the limit(in bps or pps), the later Packets will be put in a buffer
// Once the buffer overflow, later Packets will be discarded
// The the buffer limit will not be accurate, usually a little lower than specified, since it takes time(usually less than microseconds) to send Packets
type Restrict struct {
	BasicNode
	ppsLimit, bpsLimit                float64
	bufferSizeLimit, bufferCountLimit uint64
	bufferSize, bufferCount           *atomic.Uint64
	emitTime                          time.Time
}

// NewRestrict create a new restrict with the given parameter
// Next, recordSize, onEmitCallback the same as BasicNode
// ppsLimit, bpsLimit: the limit of Packets per second/bytes per second
// bufferSizeLimit, bufferCountLimit: the limit of waiting Packets, in bytes/Packets
func NewRestrict(recordSize int, onEmitCallback OnEmitCallback,
	ppsLimit, bpsLimit float64,
	bufferSizeLimit, bufferCountLimit uint64) *Restrict {
	return &Restrict{
		BasicNode:        *NewBasicNode(recordSize, onEmitCallback),
		ppsLimit:         ppsLimit,
		bpsLimit:         bpsLimit,
		bufferSizeLimit:  bufferSizeLimit,
		bufferCountLimit: bufferCountLimit,
		bufferSize:       atomic.NewUint64(0),
		bufferCount:      atomic.NewUint64(0),
		emitTime:         Now(),
	}
}

func (r *Restrict) Emit(packet *SimulatedPacket) {
	r.BasicNode.Emit(packet)
	r.bufferSize.Sub(uint64(len(packet.Actual.Data)))
	r.bufferCount.Dec()
}

func (r *Restrict) Send(packet *Packet) {
	if r.bufferSize.Load() >= r.bufferSizeLimit || r.bufferCount.Load() >= r.bufferCountLimit {
		return
	}
	sentTime := Now()
	if r.emitTime.Before(sentTime) {
		r.emitTime = sentTime
	}
	emitTime := r.emitTime
	p := &SimulatedPacket{Actual: packet, EmitTime: emitTime, SentTime: sentTime, Loss: false, Where: r}
	r.buffer.Insert(p)
	step := math.Max(1.0/r.ppsLimit, float64(len(packet.Data))/r.bpsLimit)
	r.emitTime = emitTime.Add(time.Duration(step*1000*1000) * time.Microsecond)
	r.bufferSize.Add(uint64(len(packet.Data)))
	r.bufferCount.Inc()
	r.BasicNode.OnSend(p)
}
