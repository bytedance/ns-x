package byte_ns

const preserveBufferSize = 10

// Endpoint is a node to receive Packets, Packets reached an endpoint will no longer be transmitted
type Endpoint struct {
	BasicNode
	concurrentBuffer *PacketBuffer
	localBuffer      []*SimulatedPacket
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		BasicNode:        BasicNode{},
		concurrentBuffer: NewPacketBuffer(),
		localBuffer:      make([]*SimulatedPacket, 0, preserveBufferSize),
	}
}

func (e *Endpoint) Send(packet *Packet) {
	t := Now()
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
