package networksimulator

import (
	"math"
	"time"
)

// SpeedController 模拟实际网卡收发能力的上限
type SpeedController struct {
	BasicNode
	ppsLimit float64
	bpsLimit float64
	emitTime time.Time
}

func NewSpeedController(next Node, recordSize int, onEmitCallback OnEmitCallback, ppsLimit, bpsLimit float64) *SpeedController {
	return &SpeedController{
		BasicNode: *NewBasicNode(next, recordSize, onEmitCallback),
		ppsLimit:  ppsLimit,
		bpsLimit:  bpsLimit,
		emitTime:  Now(),
	}
}

func (s *SpeedController) Send(packet *Packet) {
	sentTime := Now()
	if s.emitTime.Before(sentTime) {
		s.emitTime = sentTime
	}
	emitTime := s.emitTime
	p := &SimulatedPacket{Actual: packet, EmitTime: emitTime, SentTime: sentTime, Loss: false, Where: s}
	s.buffer.Insert(p)
	step := math.Max(1.0/s.ppsLimit, float64(len(packet.data))/s.bpsLimit)
	s.emitTime = emitTime.Add(time.Duration(step*1000*1000) * time.Microsecond)
	s.BasicNode.Send(p)
}
