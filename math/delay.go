package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math/rand"
	"time"
)

// Delay in the network
type Delay interface {
	Delay() time.Duration   // Actual Delay duration (include jitter)
	Average() time.Duration // Average delay duration
	Jitter() time.Duration  // Jitter delay duration
}

// basicDelay skeleton implementation of Delay
type basicDelay struct {
	average time.Duration
	jitter  time.Duration
	random  *rand.Rand
}

func (bd *basicDelay) Average() time.Duration {
	return bd.average
}

func (bd *basicDelay) Jitter() time.Duration {
	jitter := float64(bd.jitter.Microseconds())
	return time.Duration((2*bd.random.Float64()-1)*jitter) * time.Microsecond
}

// FixedDelay delay fixed duration
type FixedDelay struct {
	*basicDelay
}

var _ Delay = &FixedDelay{}

func NewFixedDelay(average time.Duration, jitter time.Duration, random *rand.Rand) node.PacketHandler {
	delay := &FixedDelay{
		&basicDelay{
			average: average,
			jitter:  jitter,
			random:  random,
		},
	}
	return delay.PacketHandler
}

func (nd *FixedDelay) Delay() time.Duration {
	return nd.Average() + nd.Jitter()
}

func (nd *FixedDelay) PacketHandler(base.Packet) (time.Duration, bool) {
	return nd.Delay(), false
}

// NormalDelay delay duration of normal distribution
type NormalDelay struct {
	*basicDelay
	sigma time.Duration
}

var _ Delay = &NormalDelay{}

func NewNormalDelay(average, jitter, sigma time.Duration, random *rand.Rand) node.PacketHandler {
	delay := &NormalDelay{
		&basicDelay{
			average: average,
			jitter:  jitter,
			random:  random,
		},
		sigma,
	}
	return delay.PacketHandler
}

func (nd *NormalDelay) Delay() time.Duration {
	sigma := float64(nd.sigma.Microseconds())
	return time.Duration(nd.random.NormFloat64()*sigma)*time.Microsecond + nd.Average() + nd.Jitter()
}

func (nd *NormalDelay) PacketHandler(base.Packet) (time.Duration, bool) {
	return nd.Delay(), false
}

// UniformDelay delay duration of uniform distribution
type UniformDelay struct {
	*basicDelay
}

var _ Delay = &UniformDelay{}

func NewUniformDelay(average, jitter time.Duration, random *rand.Rand) node.PacketHandler {
	delay := &UniformDelay{
		&basicDelay{
			average: average,
			jitter:  jitter,
			random:  random,
		},
	}
	return delay.PacketHandler
}

func (ud *UniformDelay) Delay() time.Duration {
	average := float64(ud.average.Microseconds())
	return time.Duration(ud.random.Float64()*average*2)*time.Microsecond + ud.Jitter()
}

func (ud *UniformDelay) PacketHandler(base.Packet) (time.Duration, bool) {
	return ud.Delay(), false
}
