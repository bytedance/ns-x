package base

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

const TestBucketSize = time.Second
const TestBucketsLimit = 128

func TestEventQueue(t *testing.T) {
	now := time.Now()
	eventQueue := NewEventQueue(TestBucketSize, TestBucketsLimit)
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

func BenchmarkEventQueue(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(TestBucketSize, TestBucketsLimit)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
	}
	for !eventQueue.IsEmpty() {
		_ = eventQueue.Dequeue()
	}
}

func BenchmarkEnqueue(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(TestBucketSize, TestBucketsLimit)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
	}
}

func BenchmarkDequeue(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(TestBucketSize, TestBucketsLimit)
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

func BenchmarkEnqueueAndDequeue(b *testing.B) {
	now := time.Now()
	eventQueue := NewEventQueue(TestBucketSize, TestBucketsLimit)
	for i := 0; i < b.N; i++ {
		eventQueue.Enqueue(NewFixedEvent(func(t time.Time) []Event {
			return nil
		}, now.Add(time.Duration(rand.Int()%20)*time.Second)))
		_ = eventQueue.Dequeue()
	}
}
