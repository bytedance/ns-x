package base

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	now := time.Now()
	eventQueue := NewEventQueue(time.Second)
	count := 50
	for i := 0; i < count; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
	}
	assert.Equal(t, count, eventQueue.Length())
	last := int64(math.MinInt64)
	for !eventQueue.IsEmpty() {
		p := eventQueue.Dequeue()
		unix := p.Time().Unix()
		assert.GreaterOrEqual(t, unix, last)
		last = unix
	}
}

func BenchmarkHeap(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(time.Second)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
	}
	for !eventQueue.IsEmpty() {
		_ = eventQueue.Dequeue()
	}
}

func BenchmarkPush(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(time.Second)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
	}
}

func BenchmarkPop(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(time.Second)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
	}
	b.ResetTimer()
	for !eventQueue.IsEmpty() {
		_ = eventQueue.Dequeue()
	}
}

func BenchmarkPushAndPop(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(time.Second)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
		_ = eventQueue.Dequeue()
	}
}

func BenchmarkLimitedHeap(b *testing.B) {
	now := time.Now()
	N := b.N
	limit := 100
	eventQueue := NewEventQueue(time.Second)
	for N > 0 {
		if N < limit {
			limit = N
		}
		for i := 0; i < limit; i++ {
			eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
				return nil
			}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
		}
		for !eventQueue.IsEmpty() {
			_ = eventQueue.Dequeue()
		}
		N -= limit
	}
}
