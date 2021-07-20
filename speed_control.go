package networksimulator

import (
	"math"
	"time"
)

type SpeedCalculator func(p float64, i float64, d float64) float64

// SpeedController 模拟实际网卡收发能力的上限
type SpeedController struct {
	*BasicNode
	ppsLimit float64
	bpsLimit float64
	emitTime time.Time
}

func (s *SpeedController) Send(packet *Packet) {
	sentTime := time.Now()
	if s.emitTime.Before(sentTime) {
		s.emitTime = sentTime
	}
	emitTime := s.emitTime
	p := &SimulatedPacket{actual: packet, emitTime: emitTime, sentTime: sentTime, loss: false, node: s}
	s.buffer.Insert(p)
	step := math.Max(1.0/s.ppsLimit, float64(len(packet.data))/s.bpsLimit)
	s.emitTime = emitTime.Add(time.Duration(step*1000*1000) * time.Microsecond)
	s.BasicNode.Send(p)
}
