package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math/rand"
	"time"
)

// some commonly used delay model

// NewFixedDelay always delay given duration
func NewFixedDelay(delay time.Duration) node.Delay {
	return func(base.Packet) time.Duration {
		return delay
	}
}

// NewNormalDelay delay with a normal distribution
func NewNormalDelay(average, sigma time.Duration, random *rand.Rand) node.Delay {
	return func(packet base.Packet) time.Duration {
		return average + time.Duration(random.NormFloat64()*float64(sigma))
	}
}

// NewUniformDelay delay with a uniform distribution
func NewUniformDelay(average time.Duration, random *rand.Rand) node.Delay {
	return func(packet base.Packet) time.Duration {
		return time.Duration(random.Int63n(int64(2 * average)))
	}
}
