package node

import (
	"byte-ns/base"
	time2 "byte-ns/time"
	"time"
)

// PacketHandler handles how much a Packet delayed and whether lost according to historical records
type PacketHandler func(packet *base.Packet, record *base.PacketQueue) (delay time.Duration, lost bool)

// Combine given handlers to sum up all their delays and losses
func Combine(handlers ...PacketHandler) PacketHandler {
	return func(packet *base.Packet, record *base.PacketQueue) (time.Duration, bool) {
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

// ChannelNode is a simulated network channel with loss, delay and reorder features
type ChannelNode struct {
	BasicNode
	handler PacketHandler
}

// NewChannelNode creates a new channel
func NewChannelNode(name string, recordSize int, onEmitCallback base.OnEmitCallback, handler PacketHandler) *ChannelNode {
	return &ChannelNode{
		BasicNode: *NewBasicNode(name, recordSize, onEmitCallback),
		handler:   handler,
	}
}

func (c *ChannelNode) Send(packet *base.Packet) {
	now := time2.Now()
	t := now
	l := false
	if c.handler != nil {
		delay, loss := c.handler(packet, c.record)
		t = t.Add(delay)
		l = l || loss
	}
	p := &base.SimulatedPacket{Actual: packet, EmitTime: t, SentTime: now, Loss: l, Where: c}
	c.Packets().Insert(p)
	c.OnSend(p)
}
