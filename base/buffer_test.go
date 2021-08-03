package base

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkBuffer(b *testing.B) {
	buffer := NewPacketBuffer()
	for i := 0; i < b.N; i++ {
		buffer.Insert(&SimulatedPacket{EmitTime: time.Unix(rand.Int63(), 0)})
	}
	buffer.Reduce(func(packet *SimulatedPacket) {
	})
}
