package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueue(t *testing.T) {
	length := 10
	it := 0
	ring := NewQueue(length)
	data := make([]interface{}, length)
	for i := 0; i < length; i++ {
		data[i] = i
		ring.Enqueue(data[i])
	}
	assert.Equal(t, length, ring.Length())
	for !ring.IsEmpty() {
		assert.Equal(t, data[it], ring.Dequeue())
		it++
	}
}

func TestQueueOverflow(t *testing.T) {
	ringLength := 100
	overflowCycle := 3
	length := ringLength * overflowCycle
	ring := NewQueue(ringLength)
	data := make([]interface{}, length)
	for i := 0; i < length; i++ {
		data[i] = i
		ring.Enqueue(data[i])
	}
	assert.Equal(t, length, ring.Length())
	it := 0
	for !ring.IsEmpty() {
		assert.Equal(t, data[it], ring.Dequeue())
		it++
	}
}
