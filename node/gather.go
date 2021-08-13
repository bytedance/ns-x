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

func (n *GatherNode) Emit(packet base.Packet, now time.Time) {
	if n.next == nil || len(n.next) != 1 {
		panic("gather node can only has single connection")
	}
	n.Events().Insert(base.NewFixedEvent(func() {
		n.ActualEmit(packet, n.next[0], now)
	}, now))
}
