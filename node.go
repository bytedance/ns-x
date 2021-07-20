package networksimulator

// Node Indicates a simulated node in the network
type Node interface {
	Send(packet *Packet) // Send a packet to the node
	packets() *PacketBuffer
	emit(packet *SimulatedPacket)
}

// OnEmitCallback called when a packet is emitted
type OnEmitCallback func(packet *SimulatedPacket)

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

func (n *BasicNode) Send(packet *SimulatedPacket) {
	n.record.Enqueue(packet)
}

func (n *BasicNode) packets() *PacketBuffer {
	return n.buffer
}

func (n *BasicNode) emit(packet *SimulatedPacket) {
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
