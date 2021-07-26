package networksimulator

// Node Indicates a simulated node in the network
type Node interface {
	Send(packet *Packet) // Send a packet to the node
	Packets() *PacketBuffer
	Emit(packet *SimulatedPacket)
}

// OnEmitCallback called when a packet is emitted
type OnEmitCallback func(packet *SimulatedPacket)
