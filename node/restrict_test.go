package node

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRestrictNodePPSLimit(t *testing.T) {
	queueLimit := int64(10)
	node := NewRestrictNode(WithPPSLimit(1.0, queueLimit))
	node.SetNext(NewEndpointNode())
	current := time.Now()
	busy := current
	for i := int64(0); i <= queueLimit; i++ {
		node.Transfer(base.RawPacket{}, current)
		busy = busy.Add(time.Second)
		assert.Equal(t, node.BusyTime(), busy)
		assert.Equal(t, node.QueuePackets(), i)
	}
	for i := int64(0); i < queueLimit; i++ {
		node.Transfer(base.RawPacket{}, current)
		assert.Equal(t, node.BusyTime(), busy)
		assert.Equal(t, node.QueuePackets(), queueLimit)
	}
}

func TestRestrictNodeBPSLimit(t *testing.T) {
	queueLimit := int64(512)
	node := NewRestrictNode(WithBPSLimit(1.0, queueLimit))
	node.SetNext(NewEndpointNode())
	current := time.Now()
	busy := current
	for i := int64(0); i <= queueLimit; i++ {
		node.Transfer(base.RawPacket{0x01}, current)
		busy = busy.Add(time.Second)
		assert.Equal(t, node.BusyTime(), busy)
		assert.Equal(t, node.QueuePackets(), i)
	}
	for i := int64(0); i < queueLimit; i++ {
		node.Transfer(base.RawPacket{0x01}, current)
		assert.Equal(t, node.BusyTime(), busy)
		assert.Equal(t, node.QueuePackets(), queueLimit)
	}
}
