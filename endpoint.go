package networksimulator

import "time"

const preserveBufferSize = 10

type Endpoint struct {
	*BasicNode
	concurrentBuffer *PacketBuffer
	localBuffer      []*SimulatedPacket
}

func NewEndPoint() *Endpoint {
	return &Endpoint{
		BasicNode:        NewBasicNode(nil, 0, nil),
		concurrentBuffer: NewPackerBuffer(),
		localBuffer:      make([]*SimulatedPacket, 0, preserveBufferSize),
	}
}

func (e *Endpoint) Send(packet *Packet) {
	t := time.Now()
	p := &SimulatedPacket{
		Actual:   packet,
		SentTime: t,
		EmitTime: t,
		Loss:     false,
		Where:    e,
	}
	e.concurrentBuffer.Insert(p)
}

func (e *Endpoint) Receive() *SimulatedPacket {
	if len(e.localBuffer) == 0 {
		e.concurrentBuffer.Reduce(func(packet *SimulatedPacket) {
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
