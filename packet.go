package networksimulator

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// Packet Indicates an Actual packet, with its data and address
type Packet struct {
	data    []byte
	address net.Addr
}

// SimulatedPacket Indicates a simulated packet, with its Actual packet and some simulated environment
type SimulatedPacket struct {
	Actual   *Packet   // the Actual packet
	EmitTime time.Time // when this packet is emitted (Where emit a packet means the packet leaves the Where, send to the next Where)
	SentTime time.Time // when this packet is sent (Where send a packet means the packet enters the Where, waiting to emit)
	Loss     bool      // whether this packet is lost
	Where    Node      // where is the packet
}

func (packet *SimulatedPacket) String() string {
	builder := strings.Builder{}
	builder.WriteString("sent time: ")
	builder.WriteString(packet.SentTime.String())
	builder.WriteRune('\n')
	builder.WriteString("emit time: ")
	builder.WriteString(packet.EmitTime.String())
	builder.WriteRune('\n')
	builder.WriteString("loss: ")
	builder.WriteString(strconv.FormatBool(packet.Loss))
	builder.WriteRune('\n')
	return builder.String()
}
