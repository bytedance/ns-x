package base

import (
	"testing"
	"time"
)

func BenchmarkBuffer(b *testing.B) {
	buffer := NewEventBuffer()
	t := time.Now()
	callback := func(event Event) {
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Insert(NewFixedEvent(nil, t))
	}
	buffer.Reduce(callback)
}
