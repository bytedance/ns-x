package node

import (
	"byte-ns/base"
	"byte-ns/time"
)

// Gather ...
type Gather struct {
	BasicNode
}

func NewGather(name string) *Gather {
	return &Gather{BasicNode{name: name}}
}

func (g *Gather) Send(packet *base.Packet) {
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
