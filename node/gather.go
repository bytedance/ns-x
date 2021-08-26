package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// GatherNode transfer packets from multiple sources to a single target
type GatherNode struct {
	*BasicNode
}

func NewGatherNode(name string) *GatherNode {
	return &GatherNode{&BasicNode{name: name}}
}

func (n *GatherNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	if len(n.next) != 1 {
		panic("gather node can only has single connection")
	}
	return base.Aggregate(
		base.NewFixedEvent(func(t time.Time) []base.Event {
			return n.ActualTransfer(packet, n, n.next[0], t)
		}, now),
	)
}
