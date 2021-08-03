package base

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// Packet indicates an Actual packet, with its data and Address
type Packet struct {
	Data    []byte
	Address net.Addr
}

// SimulatedPacket indicates a simulated packet, with its Actual packet and some simulated environment
type SimulatedPacket struct {
	Actual   *Packet   // the Actual packet
	EmitTime time.Time // when this packet is emitted (Where OnEmit a packet means the packet leaves the Where, send to the next Where)
	SentTime time.Time // when this packet is sent (Where send a packet means the packet enters the Where, waiting to OnEmit)
	Loss     bool      // whether this packet is lost
	Where    Node      // where is the packet
}

func (packet *SimulatedPacket) String() string {
	builder := strings.Builder{}
	builder.WriteString("Sent time: ")
	builder.WriteString(packet.SentTime.String())
	builder.WriteRune('\n')
	builder.WriteString("Emit time: ")
	builder.WriteString(packet.EmitTime.String())
	builder.WriteRune('\n')
	builder.WriteString("Loss: ")
	builder.WriteString(strconv.FormatBool(packet.Loss))
	builder.WriteRune('\n')
	builder.WriteString("Where: ")
	builder.WriteString(packet.Where.Name())
	builder.WriteRune('\n')
	builder.WriteString("Target: ")
	if packet.Actual.Address != nil {
		builder.WriteString(packet.Actual.Address.String())
	}
	builder.WriteRune('\n')
	return builder.String()
}