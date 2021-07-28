package byte_ns

// BasicNode implements a basic Node
type BasicNode struct {
	next           []Node
	buffer         *PacketBuffer
	record         *PacketQueue
	onEmitCallback OnEmitCallback
}

// NewBasicNode creates a new BasicNode
func NewBasicNode(recordSize int, onEmitCallback OnEmitCallback) *BasicNode {
	return &BasicNode{
		next:           []Node{},
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

func (n *BasicNode) OnEmit(packet *SimulatedPacket) {
	if n.onEmitCallback != nil {
		n.onEmitCallback(packet)
	}
}

func (n *BasicNode) Emit(packet *SimulatedPacket) {
	n.OnEmit(packet)
	if packet.Loss {
		return
	}
	for _, node := range n.next {
		node.Send(packet.Actual)
	}

}

func (n *BasicNode) Send(packet *Packet) {
	panic("not implemented")
}

func (n *BasicNode) GetNext() []Node {
	return n.next
}

func (n *BasicNode) SetNext(nodes ...Node) {
	if nodes == nil {
		panic("cannot set next nodes to nil")
	}
	n.next = nodes
}
