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

// NewChannelNode creates a new ChannelNode with the given options
func NewChannelNode(options ...Option) *ChannelNode {
	n := &ChannelNode{
		BasicNode: &BasicNode{},
	}
	apply(n, options...)
	return n
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
				return n.actualTransfer(packet, n, n.GetNext()[0], t)
			}, delay, now),
		)
	}
	return nil
}

func (n *ChannelNode) Check() {
	if len(n.GetNext()) != 1 {
		panic("channel node can only has single connection")
	}
	n.BasicNode.Check()
}

// WithPacketHandler create an option to set/overwrite the given packet handler to nodes applied
// node applied must be a ChannelNode
func WithPacketHandler(handler PacketHandler) Option {
	return func(node base.Node) {
		n, ok := node.(*ChannelNode)
		if !ok {
			panic("cannot set packet handler")
		}
		n.handler = handler
	}
}
