package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math/rand"
	"time"
)

// some commonly used reorder model

// NewNormalReorder for correlation possibility, reorder same to last packet, or reorder with the given possibility
// reorder means the packet will be sent delta time in advance
func NewNormalReorder(delta time.Duration, possibility, correlation float64, random *rand.Rand) node.Reorder {
	last := false
	return func(base.Packet) time.Duration {
		if random.Float64() >= correlation {
			last = random.Float64() < possibility
		}
		if last {
			return -delta
		}
		return 0
	}
}

// NewGapReorder for following gap packets after a reorder packet, no reorder; otherwise same to normal reorder
func NewGapReorder(delta time.Duration, possibility, correlation float64, gap int, random *rand.Rand) node.Reorder {
	count := 0
	last := false
	return func(base.Packet) time.Duration {
		count++
		if count < gap {
			last = false
			return 0
		}
		if random.Float64() >= correlation {
			last = random.Float64() < possibility
		}
		if last {
			count = 0
			return -delta
		}
		return 0
	}
}
