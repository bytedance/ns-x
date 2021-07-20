package networksimulator

import (
	"time"
)

// PacketHandler 定义了网损模块
// 返回值：delayTime: 经过该模块的延迟; isLoss: 是否丢包, True为丢包
type PacketHandler func(packet *Packet, record *PacketQueue) (delayTime time.Duration, isLoss bool)

// Channel indicates a simulated channel, with loss, delay and reorder
type Channel struct {
	*BasicNode
	handlers []PacketHandler
}

func NewChannel(next Node, recordSize int, onEmitCallback OnEmitCallback, handlers []PacketHandler) *Channel {
	return &Channel{
		BasicNode: NewBasicNode(next, recordSize, onEmitCallback),
		handlers:  handlers,
	}
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
	p := &SimulatedPacket{Actual: packet, EmitTime: t, SentTime: now, Loss: loss}
	c.buffer.Insert(p)
	c.BasicNode.Send(p)
}
