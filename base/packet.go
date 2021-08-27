package base

// Packet in the network
type Packet interface {
	Size() int
}

type RawPacket []byte

func (p RawPacket) Size() int {
	return len(p)
}

type SimulatePacket struct {
	Data   Packet
	Source Node
	Target Node
}

func (p *SimulatePacket) Size() int {
	return p.Data.Size()
}

// IPPacket is a packet for the ip protocol (ipv4)
// DO NOT USE NOW, need further support
type IPPacket struct {
	Version            byte   // 4bit
	HeaderSize         byte   // 4 bit
	ServiceType        byte   // 8 bit
	TotalSize          uint16 // 16 bit
	Identifier         uint16 // 16 bit
	Flags              byte   // 3 bit
	FragmentOffset     uint16 // 13 bit
	TTL                byte   // 8 bit
	Protocol           byte   // 8 bit
	HeaderChecksum     uint16 // 16 bit
	SourceAddress      uint32 // 32 bit
	DestinationAddress uint32 // 32 bit
	Options            []byte // vary
	Data               Packet
}

func (p *IPPacket) Size() int {
	return int(p.TotalSize)
}

// UDPPacket is a packet for the udp protocol
// DO NOT USE NOW, need further support
type UDPPacket struct {
	// Size field has no use
	SourcePort uint16 // 2 byte
	TargetPort uint16 // 2 byte
	Checksum   uint16 // 2 byte
	Data       Packet
}

func (p *UDPPacket) Size() int {
	return p.Data.Size() + 8
}
