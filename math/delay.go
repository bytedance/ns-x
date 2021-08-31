package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math"
	"math/rand"
	"time"
)

// some commonly used delay model

// NewFixedDelay always delay given duration
func NewFixedDelay(delay time.Duration) node.Delay {
	if delay < 0 {
		panic("invalid argument")
	}
	return func(base.Packet) time.Duration {
		return delay
	}
}

// NewNormalDelay delay with a normal distribution
func NewNormalDelay(average, sigma time.Duration, random *rand.Rand) node.Delay {
	if sigma < 0 || random == nil {
		panic("invalid argument")
	}
	return func(base.Packet) time.Duration {
		return average + time.Duration(random.NormFloat64()*float64(sigma))
	}
}

// NewUniformDelay delay with a uniform distribution in [0, 2*average)
func NewUniformDelay(average time.Duration, random *rand.Rand) node.Delay {
	if average <= 0 || random == nil {
		panic("invalid argument")
	}
	return func(base.Packet) time.Duration {
		return time.Duration(random.Int63n(int64(2 * average)))
	}
}

// NewParetoDelay delay with a pareto distribution, see https://en.wikipedia.org/wiki/Pareto_distribution
func NewParetoDelay(minDelay time.Duration, alpha float64, random *rand.Rand) node.Delay {
	if minDelay <= 0 || alpha <= 0 || random == nil {
		panic("invalid argument")
	}
	return func(base.Packet) time.Duration {
		return time.Duration(float64(minDelay) * math.Pow(random.Float64(), -1/alpha))
	}
}
