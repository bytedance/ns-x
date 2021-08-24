package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math/rand"
	"time"
)

// Reorder change the delay time to reorder packets
type Reorder interface {
	Reorder() time.Duration // delta delay duration
}

// NoneReorder no reorder
type NoneReorder struct {
}

var _ Reorder = &NoneReorder{}

func NewNoneReorder() node.PacketHandler {
	reorder := &NoneReorder{}
	return reorder.PacketHandler
}

func (nr *NoneReorder) Reorder() time.Duration {
	return 0
}

func (nr *NoneReorder) PacketHandler(base.Packet) (time.Duration, bool) {
	return nr.Reorder(), false
}

// NormalReorder reorder state will be same as the last packet with correlation possibility
// otherwise reorder with the given possibility
type NormalReorder struct {
	delay       time.Duration
	possibility float64
	correlation float64
	lastReorder bool
	random      *rand.Rand
}

var _ Reorder = &NormalReorder{}

func NewNormalReorder(delay time.Duration, possibility, correlation float64, random *rand.Rand) node.PacketHandler {
	reorder := &NormalReorder{
		delay:       delay,
		possibility: possibility,
		correlation: correlation,
		lastReorder: false, // packets before the first packet are regarded as not reordered
		random:      random,
	}
	return reorder.PacketHandler
}

func (nr *NormalReorder) Reorder() time.Duration {
	if nr.random.Float64() < nr.correlation {
		if nr.lastReorder {
			return -nr.delay
		} else {
			return 0
		}
	}
	nr.lastReorder = nr.random.Float64() < nr.possibility
	if nr.lastReorder {
		return -nr.delay
	}
	return 0
}

func (nr *NormalReorder) PacketHandler(base.Packet) (time.Duration, bool) {
	return nr.Reorder(), false
}

// GapReorder gap reorder model, similar to tc
type GapReorder struct {
	delay       time.Duration
	possibility float64
	correlation float64
	lastReorder bool
	gap         int
	count       int
	random      *rand.Rand
}

var _ Reorder = &GapReorder{}

func NewGapReorder(delay time.Duration, possibility,
	correlation float64, gap int, random *rand.Rand) node.PacketHandler {
	reorder := &GapReorder{
		delay:       delay,
		possibility: possibility,
		correlation: correlation,
		lastReorder: false, // packets before the first packet are regarded as not reordered
		gap:         gap,
		count:       0,
		random:      random,
	}
	return reorder.PacketHandler
}

func (gr *GapReorder) Reorder() time.Duration {
	if gr.random.Float64() < gr.correlation {
		if gr.lastReorder {
			return -gr.delay
		}
		return 0
	}
	gr.count++
	if gr.count < gr.gap {
		gr.lastReorder = false
		return 0
	}
	gr.lastReorder = gr.random.Float64() < gr.possibility
	if gr.lastReorder {
		gr.count = 0
		return -gr.delay
	}
	return 0
}

func (gr *GapReorder) PacketHandler(base.Packet) (time.Duration, bool) {
	return gr.Reorder(), false
}
