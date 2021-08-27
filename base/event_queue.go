package base

import (
	"container/heap"
	"time"
)

// EventQueue is used to sort events according to the time of events
type EventQueue struct {
	total         int
	buckets       *Queue
	threshold     time.Time
	bucketSize    time.Duration
	maxBuckets    int
	currentBucket *bucket
	defaultBucket *bucket
}

func NewEventQueue(bucketSize time.Duration, maxBuckets int) *EventQueue {
	return &EventQueue{
		total:         0,
		buckets:       NewQueue(0),
		bucketSize:    bucketSize,
		maxBuckets:    maxBuckets,
		currentBucket: &bucket{},
		defaultBucket: &bucket{},
	}
}

func (q *EventQueue) Enqueue(event Event) {
	t := event.Time()
	if q.total <= 0 {
		q.threshold = t.Add(q.bucketSize)
	}
	b := q.currentBucket
	if t.After(q.threshold) {
		index := int(t.Sub(q.threshold) / q.bucketSize)
		if index > q.maxBuckets {
			b = q.defaultBucket
		} else {
			for index >= q.buckets.Length() {
				q.buckets.Enqueue(&bucket{})
			}
			b = q.buckets.At(index).(*bucket)
		}
	}
	heap.Push(b, event)
	q.total++
}

func (q *EventQueue) Dequeue() Event {
	if q.currentBucket.IsEmpty() {
		panic("no more events")
	}
	event := heap.Pop(q.currentBucket).(Event)
	q.total--
	for q.currentBucket.IsEmpty() {
		if q.buckets.IsEmpty() {
			break
		}
		q.currentBucket = q.buckets.Dequeue().(*bucket)
		q.threshold = q.threshold.Add(q.bucketSize)
	}
	t := q.threshold.Add(q.bucketSize * time.Duration(q.maxBuckets))
	for !q.defaultBucket.IsEmpty() {
		e := q.defaultBucket.Peek()
		if e.Time().After(t) {
			break
		}
		heap.Pop(q.defaultBucket)
		q.total--
		q.Enqueue(e)
	}
	return event
}

func (q *EventQueue) Peek() Event {
	if q.currentBucket.IsEmpty() {
		panic("no more events")
	}
	return q.currentBucket.Peek()
}

func (q *EventQueue) Length() int {
	return q.total
}

func (q *EventQueue) IsEmpty() bool {
	return q.Length() <= 0
}

type bucket struct {
	storage []Event
}

func (b *bucket) IsEmpty() bool {
	return b.Len() == 0
}

func (b *bucket) Less(i, j int) bool {
	ti := b.storage[i].Time()
	tj := b.storage[j].Time()
	return ti.Before(tj)
}

func (b *bucket) Len() int {
	return len(b.storage)
}

func (b *bucket) Swap(i, j int) {
	b.storage[i], b.storage[j] = b.storage[j], b.storage[i]
}

func (b *bucket) Push(x interface{}) {
	b.storage = append(b.storage, x.(Event))
}

func (b *bucket) Pop() interface{} {
	x := b.storage[b.Len()-1]
	b.storage = b.storage[:b.Len()-1]
	return x
}

func (b *bucket) Peek() Event {
	return b.storage[0]
}
