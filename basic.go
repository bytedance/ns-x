package networksimulator

// BasicNode implements a basic Node
type BasicNode struct {
	Next           Node
	buffer         *PacketBuffer
	record         *PacketQueue
	onEmitCallback OnEmitCallback
}

// NewBasicNode creates a new BasicNode
func NewBasicNode(recordSize int, onEmitCallback OnEmitCallback) *BasicNode {
	return &BasicNode{
		buffer:         NewPacketBuffer(),
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
	if n.Next != nil {
		n.Next.Send(packet.Actual)
	}
}

func (n *BasicNode) Send(packet *Packet) {
	panic("not implemented")
}
