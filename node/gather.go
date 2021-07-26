package node

import (
	"network-simulator/core"
)

type Gather struct {
	BasicNode
}

func (g *Gather) Send(packet *core.Packet) {
	t := core.Now()
	p := &core.SimulatedPacket{
		Actual:   packet,
		SentTime: t,
		EmitTime: t,
		Where:    g,
		Loss:     false,
	}
	g.buffer.Insert(p)
	g.BasicNode.OnSend(p)
}
