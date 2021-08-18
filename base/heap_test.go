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
	eventHeap := &EventHeap{}
	count := 50
	for i := 0; i < count; i++ {
		heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, time.Unix(rand.Int63(), 0)))
	}
	assert.Equal(t, count, eventHeap.Len())
	last := int64(math.MinInt64)
	for !eventHeap.IsEmpty() {
		p := heap.Pop(eventHeap).(Event)
		unix := p.Time().Unix()
		assert.GreaterOrEqual(t, unix, last)
		last = unix
	}
}

func BenchmarkHeap(b *testing.B) {
	eventHeap := &EventHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, time.Unix(rand.Int63(), 0)))
	}
	for !eventHeap.IsEmpty() {
		_ = heap.Pop(eventHeap).(Event)
	}
}

func BenchmarkSinglePush(b *testing.B) {
	eventHeap := &EventHeap{}
	for i := 0; i < b.N; i++ {
		eventHeap.Storage = eventHeap.Storage[0:0]
		heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, time.Unix(rand.Int63(), 0)))
	}
}

func BenchmarkPush(b *testing.B) {
	eventHeap := &EventHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, time.Unix(rand.Int63(), 0)))
	}
}

func BenchmarkPop(b *testing.B) {
	eventHeap := &EventHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, time.Unix(rand.Int63(), 0)))
	}
	b.ResetTimer()
	for !eventHeap.IsEmpty() {
		_ = heap.Pop(eventHeap).(Event)
	}
}

func BenchmarkPushAndPop(b *testing.B) {
	eventHeap := &EventHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, time.Unix(rand.Int63(), 0)))
		_ = heap.Pop(eventHeap).(Event)
	}
}

func BenchmarkLimitedHeap(b *testing.B) {
	N := b.N
	limit := 100
	eventHeap := &EventHeap{}
	for N > 0 {
		if N < limit {
			limit = N
		}
		for i := 0; i < limit; i++ {
			heap.Push(eventHeap, NewFixedEvent(func(t time.Time) []Event {
				return nil
			}, time.Unix(rand.Int63(), 0)))
		}
		for !eventHeap.IsEmpty() {
			_ = heap.Pop(eventHeap).(Event)
		}
		N -= limit
	}
}
