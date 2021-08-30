package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

// GatherNode transfer packets from multiple sources to a single target
type GatherNode struct {
	*BasicNode
}

// NewGatherNode create a gather node
func NewGatherNode(options ...Option) *GatherNode {
	n := &GatherNode{&BasicNode{}}
	apply(n, options...)
	return n
}

func (n *GatherNode) Transfer(packet base.Packet, now time.Time) []base.Event {
	if len(n.GetNext()) != 1 {
		panic("gather node can only has single connection")
	}
	return base.Aggregate(
		base.NewFixedEvent(func(t time.Time) []base.Event {
			return n.actualTransfer(packet, n, n.GetNext()[0], t)
		}, now),
	)
}
