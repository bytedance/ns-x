package networksimulator

import (
	"time"
)

// PacketHandler 定义了网损模块
// 返回值：delayTime: 经过该模块的延迟; isLoss: 是否丢包, True为丢包
type PacketHandler func(packet *Packet, record *PacketQueue) (delayTime time.Duration, isLoss bool)

type Channel struct {
	*BasicNode
	handlers []PacketHandler
}

func NewChannel() *Channel {
	return &Channel{}
}

func (c *Channel) Send(packet *Packet) {
	now := time.Now()
	t := now
	loss := false
	for _, h := range c.handlers {
		duration, l := h(packet, c.record)
		if l {
			loss = true
		}
		t = t.Add(duration)
	}
	p := &SimulatedPacket{actual: packet, emitTime: t, sentTime: now, loss: loss}
	c.buffer.Insert(p)
	c.BasicNode.Send(p)
}
