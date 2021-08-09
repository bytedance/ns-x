package node

import (
	"byte-ns/base"
	"byte-ns/time"
)

const preserveBufferSize = 10

// EndpointNode is a node to receive Packets, Packets reached an endpoint will no longer be transmitted
type EndpointNode struct {
	BasicNode
	concurrentBuffer *base.PacketBuffer
	localBuffer      []*base.SimulatedPacket
}

func NewEndpointNode(name string) *EndpointNode {
	return &EndpointNode{
		BasicNode:        BasicNode{name: name},
		concurrentBuffer: base.NewPacketBuffer(),
		localBuffer:      make([]*base.SimulatedPacket, 0, preserveBufferSize),
	}
}

func (e *EndpointNode) Send(packet []byte) {
	t := time.Now()
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
func (e *EndpointNode) Receive() *base.SimulatedPacket {
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
