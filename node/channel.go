package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// PacketHandler handles how much a Packet delayed and whether lost according to historical records
type PacketHandler func(packet base.Packet) (delay time.Duration, lost bool)

// Combine given handlers to sum up all their delays and losses
func Combine(handlers ...PacketHandler) PacketHandler {
	return func(packet base.Packet) (time.Duration, bool) {
		delay := time.Duration(0)
		loss := false
		for _, handler := range handlers {
			d, l := handler(packet)
			delay += d
			loss = l || loss
		}
		return delay, loss
	}
}

// ChannelNode is a simulated network channel with loss, delay and reorder features
type ChannelNode struct {
	*BasicNode
	handler PacketHandler
}

// NewChannelNode creates a new channel
func NewChannelNode(name string, onEmitCallback base.TransferCallback, handler PacketHandler) *ChannelNode {
	return &ChannelNode{
		BasicNode: NewBasicNode(name, onEmitCallback),
		handler:   handler,
	}
}

func (n *ChannelNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	delay := time.Duration(0)
	loss := false
	if n.handler != nil {
		d, l := n.handler(packet)
		delay += d
		loss = loss || l
	}
	if !loss {
		return base.Aggregate(
			base.NewDelayedEvent(func(t time.Time) []base.Event {
				return n.ActualTransfer(packet, n.next[0], t)
			}, delay, now),
		)
	}
	return nil
}

func (n *ChannelNode) Check() {
	if len(n.next) != 1 {
		panic("channel node can only has single connection")
	}
	n.BasicNode.Check()
}
