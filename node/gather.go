package node

import (
	"ns-x/base"
	"time"
)

// GatherNode ...
type GatherNode struct {
	*BasicNode
}

func NewGatherNode(name string) *GatherNode {
	return &GatherNode{&BasicNode{name: name}}
}

func (n *GatherNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	if n.next == nil || len(n.next) != 1 {
		panic("gather node can only has single connection")
	}
	return base.Aggregate(
		base.NewFixedEvent(func() []base.Event {
			return n.ActualEmit(packet, n.next[0], now)
		}, now),
	)
}
