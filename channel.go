package networksimulator

import (
	"time"
)

// PacketHandler handles how much a Packet delayed and whether lost in a PacketQueue
type PacketHandler func(packet *Packet, record *PacketQueue) (delay time.Duration, lost bool)

// Combine given handlers to sum up all their delays and losses
func Combine(handlers ...PacketHandler) PacketHandler {
	return func(packet *Packet, record *PacketQueue) (time.Duration, bool) {
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

// Channel is a simulated network channel with loss, delay and reorder features
type Channel struct {
	BasicNode
	handler PacketHandler
}

// NewChannel creates a new channel
func NewChannel(next Node, recordSize int, onEmitCallback OnEmitCallback, handler PacketHandler) *Channel {
	return &Channel{
		BasicNode: *NewBasicNode(next, recordSize, onEmitCallback),
		handler:   handler,
	}
}

func (c *Channel) Send(packet *Packet) {
	now := Now()
	t := now
	l := false
	if c.handler != nil {
		delay, loss := c.handler(packet, c.record)
		t = t.Add(delay)
		l = l || loss
	}
	p := &SimulatedPacket{Actual: packet, EmitTime: t, SentTime: now, Loss: l, Where: c}
	c.buffer.Insert(p)
	c.BasicNode.OnSend(p)
}
