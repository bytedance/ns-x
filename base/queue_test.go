package base

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestQueue(t *testing.T) {
	length := 1000
	it := 0
	ring := NewDataQueue(length)
	data := make([]interface{}, length)
	for i := 0; i < length; i++ {
		data[i] = i
		ring.Enqueue(data[i])
	}
	assert.Equal(t, length, ring.Length())
	for ring.Length() > 0 {
		assert.Equal(t, data[it], ring.Dequeue())
		it++
	}
}

func TestQueueOverflow(t *testing.T) {
	ringLength := 100
	overflowCycle := 3
	length := 35
	it := ringLength*overflowCycle + length - ringLength
	ring := NewDataQueue(ringLength)
	data := make([]interface{}, ringLength*overflowCycle+length)
	for i := 0; i < ringLength*overflowCycle+length; i++ {
		data[i] = i
		ring.Enqueue(data[i])
	}
	assert.Equal(t, ringLength, ring.Length())
	for ring.Length() > 0 {
		assert.Equal(t, data[it], ring.Dequeue())
		it++
	}
}

func TestQueueOverflowShouldFail(t *testing.T) {
	ringLength := 100
	overflowCycle := 3
	length := 35
	it := ringLength*overflowCycle + length - ringLength
	ring := NewDataQueue(ringLength)
	data := make([]interface{}, ringLength*overflowCycle+length)
	for i := 0; i < ringLength*overflowCycle+length; i++ {
		data[i] = i
		ring.Enqueue(data[i])
	}
	assert.Equal(t, ringLength, ring.Length())
	for ring.Length() > 0 {
		assert.Equal(t, data[it], ring.Dequeue())
		it++
	}
}

func TestQueueOverflowWithDequeue(t *testing.T) {
	ringLength := 100
	overflowCycle := 1
	length := 37
	it := int(math.Max(float64(ringLength*overflowCycle+length-ringLength), float64(length)))
	ring := NewDataQueue(ringLength)
	data := make([]interface{}, ringLength*overflowCycle+length)
	for i := 0; i < ringLength*overflowCycle+length; i++ {
		data[i] = i
		ring.Enqueue(data[i])
	}
	assert.Equal(t, ringLength, ring.Length())
	for ring.Length() > 0 {
		assert.Equal(t, data[it], ring.Dequeue())
		it++
	}
}
