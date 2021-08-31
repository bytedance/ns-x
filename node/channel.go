package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// Loss whether the packet loss
type Loss func(packet base.Packet) bool

// Delay how long for the packet
type Delay func(packet base.Packet) time.Duration

// Reorder advance how long for the packet
type Reorder func(packet base.Packet) time.Duration

// handler handles how much a Packet delayed and whether lost
type handler func(packet base.Packet) (delay time.Duration, lost bool)

// combine given handlers to sum up all their delays and losses
func combine(handlers ...handler) handler {
	return func(packet base.Packet) (time.Duration, bool) {
		delay := time.Duration(0)
		loss := false
		for _, handler := range handlers {
			if handler == nil {
				continue
			}
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
	handler handler
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
		delay, loss = n.handler(packet)
	}
	if loss {
		return nil
	}
	if delay < 0 {
		delay = 0
	}
	return base.Aggregate(
		base.NewDelayedEvent(func(t time.Time) []base.Event {
			return n.actualTransfer(packet, n, n.GetNext()[0], t)
		}, delay, now),
	)
}

func (n *ChannelNode) Check() {
	if len(n.GetNext()) != 1 {
		panic("channel node can only has single connection")
	}
	n.BasicNode.Check()
}

// WithLoss create an Option to add a Loss on the ChannelNode applied
// node applied must be a ChannelNode
func WithLoss(loss Loss) Option {
	return func(node base.Node) {
		n, ok := node.(*ChannelNode)
		if !ok {
			panic("cannot set loss")
		}
		n.handler = combine(n.handler, func(packet base.Packet) (time.Duration, bool) {
			return 0, loss(packet)
		})
	}
}

// WithDelay create an Option to add a Delay on the ChannelNode applied
// node applied must be a ChannelNode
func WithDelay(delay Delay) Option {
	return func(node base.Node) {
		n, ok := node.(*ChannelNode)
		if !ok {
			panic("cannot set delay")
		}
		n.handler = combine(n.handler, func(packet base.Packet) (time.Duration, bool) {
			return delay(packet), false
		})
	}
}

// WithReorder create an Option to add a Reorder on the ChannelNode applied
// node applied must be a ChannelNode
func WithReorder(reorder Reorder) Option {
	return func(node base.Node) {
		n, ok := node.(*ChannelNode)
		if !ok {
			panic("cannot set reorder")
		}
		n.handler = combine(n.handler, func(packet base.Packet) (time.Duration, bool) {
			return reorder(packet), false
		})
	}
}
