package node

import (
	"network-simulator/core"
	"time"
)

// PacketHandler 定义了网损模块
// 返回值：delayTime: 经过该模块的延迟; isLoss: 是否丢包, True为丢包
type PacketHandler func(packet *core.Packet, record *core.PacketQueue) (time.Duration, bool)

// Combine the given handlers, the result will add up all the delay time and loss
func Combine(handlers ...PacketHandler) PacketHandler {
	return func(packet *core.Packet, record *core.PacketQueue) (time.Duration, bool) {
		delay := time.Duration(0)
		loss := false
		for _, handler := range handlers {
			d, l := handler(packet, record)
			delay += d
			loss = l || loss
		}
		return delay, loss
	}
}

// Channel indicates a simulated channel, with loss, delay and reorder
type Channel struct {
	BasicNode
	handler PacketHandler
}

// NewChannel creates a new channel
func NewChannel(next core.Node, recordSize int, onEmitCallback core.OnEmitCallback, handler PacketHandler) *Channel {
	return &Channel{
		BasicNode: *NewBasicNode(next, recordSize, onEmitCallback),
		handler:   handler,
	}
}

func (c *Channel) Send(packet *core.Packet) {
	now := core.Now()
	t := now
	l := false
	if c.handler != nil {
		delay, loss := c.handler(packet, c.record)
		t = t.Add(delay)
		l = l || loss
	}
	p := &core.SimulatedPacket{Actual: packet, EmitTime: t, SentTime: now, Loss: l, Where: c}
	c.buffer.Insert(p)
	c.BasicNode.OnSend(p)
}
