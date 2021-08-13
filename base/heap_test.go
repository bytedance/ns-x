package base

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	packetHeap := &EventHeap{}
	count := 50
	for i := 0; i < count; i++ {
		heap.Push(packetHeap, &SimulatedPacket{EmitTime: time.Unix(rand.Int63(), 0)})
	}
	assert.Equal(t, count, packetHeap.Len())
	last := int64(math.MinInt64)
	for !packetHeap.IsEmpty() {
		p := heap.Pop(packetHeap).(*SimulatedPacket)
		unix := p.EmitTime.Unix()
		assert.GreaterOrEqual(t, unix, last)
		last = unix
	}
}

func BenchmarkHeap(b *testing.B) {
	packetHeap := &EventHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(packetHeap, &SimulatedPacket{EmitTime: time.Unix(rand.Int63(), 0)})
	}
	for !packetHeap.IsEmpty() {
		_ = heap.Pop(packetHeap).(*SimulatedPacket)
	}
}
