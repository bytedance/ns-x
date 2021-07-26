package node

import (
	"network-simulator/core"
)

const preserveBufferSize = 10

// Endpoint is a node to receive Packets, Packets reach an endpoint will no longer transmit
type Endpoint struct {
	BasicNode
	concurrentBuffer *core.PacketBuffer
	localBuffer      []*core.SimulatedPacket
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		BasicNode:        BasicNode{},
		concurrentBuffer: core.NewPackerBuffer(),
		localBuffer:      make([]*core.SimulatedPacket, 0, preserveBufferSize),
	}
}

func (e *Endpoint) Send(packet *core.Packet) {
	t := core.Now()
	p := &core.SimulatedPacket{
		Actual:   packet,
		SentTime: t,
		EmitTime: t,
		Loss:     false,
		Where:    e,
	}
	e.concurrentBuffer.Insert(p)
}

func (e *Endpoint) Receive() *core.SimulatedPacket {
	if len(e.localBuffer) == 0 {
		e.concurrentBuffer.Reduce(func(packet *core.SimulatedPacket) {
			e.localBuffer = append(e.localBuffer, packet)
		})
	}
	index := len(e.localBuffer) - 1
	if index < 0 {
		return nil
	}
	p := e.localBuffer[index]
	e.localBuffer = e.localBuffer[:index]
	return p
}
