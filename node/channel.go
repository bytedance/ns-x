package node

import (
	"ns-x/base"
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
func NewChannelNode(name string, onEmitCallback base.OnEmitCallback, handler PacketHandler) *ChannelNode {
	return &ChannelNode{
		BasicNode: NewBasicNode(name, onEmitCallback),
		handler:   handler,
	}
}

func (n *ChannelNode) Emit(packet base.Packet, now time.Time) {
	delay := time.Duration(0)
	loss := false
	if n.handler != nil {
		d, l := n.handler(packet)
		delay += d
		loss = loss || l
	}
	if !loss {
		n.Events().Insert(base.NewDelayedEvent(func(t time.Time) {
			n.ActualEmit(packet, n.next[0], t)
		}, delay, now))
	}
}

func (n *ChannelNode) Check() {
	if n.next == nil || len(n.next) != 1 {
		panic("channel node can only has single connection")
	}
	n.BasicNode.Check()
}
