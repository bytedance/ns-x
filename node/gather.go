package node

import (
	"byte-ns/base"
	"byte-ns/time"
)

// GatherNode ...
type GatherNode struct {
	BasicNode
}

func NewGatherNode(name string) *GatherNode {
	return &GatherNode{BasicNode{name: name}}
}

func (g *GatherNode) Send(packet []byte) {
	t := time.Now()
	p := &base.SimulatedPacket{
		Actual:   packet,
		SentTime: t,
		EmitTime: t,
		Where:    g,
		Loss:     false,
	}
	g.OnSend(p)
	g.Emit(p)
}
