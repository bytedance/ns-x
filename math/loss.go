package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math/rand"
)

// some commonly used loss model

// NewRandomLoss loss with the given possibility
func NewRandomLoss(possibility float64, random *rand.Rand) node.Loss {
	return func(base.Packet) bool {
		return random.Float64() < possibility
	}
}

// NewGilbertLoss loss with gilbert-elliott model, see https://en.wikipedia.org/wiki/Burst_error
func NewGilbertLoss(g2b, b2g float64, lossG, lossB float64, random *rand.Rand) node.Loss {
	state := false // true for good state, false for bad state
	return func(base.Packet) bool {
		loss := false
		if state {
			if random.Float64() < lossG {
				loss = true
			}
			if random.Float64() < g2b {
				state = false
			}
		} else {
			if random.Float64() < lossB {
				loss = true
			}
			if random.Float64() < b2g {
				state = true
			}
		}
		return loss
	}
}
