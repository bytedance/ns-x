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
		base.NewFixedEvent(func(t time.Time) []base.Event {
			return n.ActualTransfer(packet, n.next[0], t)
		}, now),
	)
}
