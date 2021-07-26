package networksimulator

// BasicNode implement Node basically
type BasicNode struct {
	next           Node
	buffer         *PacketBuffer
	record         *PacketQueue
	onEmitCallback OnEmitCallback
}

func NewBasicNode(next Node, recordSize int, onEmitCallback OnEmitCallback) *BasicNode {
	return &BasicNode{
		next:           next,
		buffer:         NewPackerBuffer(),
		record:         NewPacketQueue(recordSize),
		onEmitCallback: onEmitCallback,
	}
}

func (n *BasicNode) OnSend(packet *SimulatedPacket) {
	n.record.Enqueue(packet)
}

func (n *BasicNode) Packets() *PacketBuffer {
	return n.buffer
}

func (n *BasicNode) Emit(packet *SimulatedPacket) {
	if n.onEmitCallback != nil {
		n.onEmitCallback(packet)
	}
	if packet.Loss {
		return
	}
	if n.next != nil {
		n.next.Send(packet.Actual)
	}
}

func (n *BasicNode) Send(packet *Packet) {
	panic("not implemented")
}
