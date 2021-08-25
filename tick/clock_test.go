package tick

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRealClock(t *testing.T) {
	clock := NewRealClock()
	N := 100
	threshold := time.Millisecond
	for i := 0; i < N; i++ {
		assert.LessOrEqual(t, time.Now().Sub(clock()), threshold)
	}
}

func TestStepClock(t *testing.T) {
	now := time.Now()
	step := time.Second
	clock := NewStepClock(now, step)
	N := 100
	for i := 0; i < N; i++ {
		assert.Equal(t, now, clock())
		now = now.Add(step)
	}
}
