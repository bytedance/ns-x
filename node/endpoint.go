package node

import (
	"byte-ns/base"
)

const preserveBufferSize = 10

// Endpoint is a node to receive Packets, Packets reached an endpoint will no longer be transmitted
type Endpoint struct {
	BasicNode
	concurrentBuffer *base.PacketBuffer
	localBuffer      []*base.SimulatedPacket
}

func NewEndpoint(name string) *Endpoint {
	return &Endpoint{
		BasicNode:        BasicNode{name: name},
		concurrentBuffer: base.NewPacketBuffer(),
		localBuffer:      make([]*base.SimulatedPacket, 0, preserveBufferSize),
	}
}

func (e *Endpoint) Send(packet *base.Packet) {
	t := base.Now()
	p := &base.SimulatedPacket{
		Actual:   packet,
		SentTime: t,
		EmitTime: t,
		Loss:     false,
		Where:    e,
	}
	e.concurrentBuffer.Insert(p)
	e.OnSend(p)
}

// Receive a packet if possible, nil otherwise
func (e *Endpoint) Receive() *base.SimulatedPacket {
	if len(e.localBuffer) == 0 {
		e.concurrentBuffer.Reduce(func(packet *base.SimulatedPacket) {
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
