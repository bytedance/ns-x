package base

import (
	"strconv"
	"strings"
	"time"
)

type Packet interface {
	Size() int
}

type RawPacket []byte

func (p RawPacket) Size() int {
	return len(p)
}

// SimulatedPacket indicates a simulated packet, with its Actual packet and some simulated environment
type SimulatedPacket struct {
	Actual   []byte    // the Actual packet
	EmitTime time.Time // when this packet is emitted (Where emit a packet means the packet leaves the Where, send to the next Where)
	SentTime time.Time // when this packet is sent (Where send a packet means the packet enters the Where, waiting to emit)
	Loss     bool      // whether this packet is lost
	Where    Node      // where is the packet
}

func (packet *SimulatedPacket) String() string {
	builder := strings.Builder{}
	builder.WriteString("Sent time: ")
	builder.WriteString(packet.SentTime.String())
	builder.WriteRune('\n')
	builder.WriteString("Transfer time: ")
	builder.WriteString(packet.EmitTime.String())
	builder.WriteRune('\n')
	builder.WriteString("Loss: ")
	builder.WriteString(strconv.FormatBool(packet.Loss))
	builder.WriteRune('\n')
	builder.WriteString("Where: ")
	builder.WriteString(packet.Where.Name())
	builder.WriteRune('\n')
	return builder.String()
}
