package byte_ns

// Node Indicates a simulated node in the network
type Node interface {
	Send(packet *Packet) // Send a packet to the node
	Packets() *PacketBuffer
	Emit(packet *SimulatedPacket)
	GetNext() []Node
	SetNext(nodes ...Node)
}

// OnEmitCallback called when a packet is emitted
type OnEmitCallback func(packet *SimulatedPacket)
