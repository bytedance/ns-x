package math

import (
	"byte-ns/base"
	node2 "byte-ns/node"
	"math/rand"
	"time"
)

// Reorder 接口，目前有不乱序，概率乱序，gap乱序三种模型
type Reorder interface {
	// Reorder 返回需要增加的延迟时间 (可能为负值)，通过将该包更"早"发送实现乱序
	Reorder() time.Duration
}

// NoneReorder 一律不乱序，不减少延迟时间
type NoneReorder struct {
}

var _ Reorder = &NoneReorder{}

func NewNoneReorder() node2.PacketHandler {
	reorder := &NoneReorder{}
	return reorder.PacketHandler
}

func (nr *NoneReorder) Reorder() time.Duration {
	return 0
}

func (nr *NoneReorder) PacketHandler(*base.Packet, *base.PacketQueue) (time.Duration, bool) {
	return nr.Reorder(), false
}

// NormalReorder 概率乱序模型
// 以 correlation 的概率和上一个包乱序情况相同
// 剩余情况以 possibility 概率乱序
type NormalReorder struct {
	delay       time.Duration
	possibility float64
	correlation float64
	lastReorder bool
	random      *rand.Rand
}

var _ Reorder = &NormalReorder{}

func NewNormalReorder(delay time.Duration, possibility, correlation float64, random *rand.Rand) node2.PacketHandler {
	reorder := &NormalReorder{
		delay:       delay,
		possibility: possibility,
		correlation: correlation,
		lastReorder: false, // 默认首包之前的包没有reorder
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

func (nr *NormalReorder) PacketHandler(*base.Packet, *base.PacketQueue) (time.Duration, bool) {
	return nr.Reorder(), false
}

// GapReorder gap 乱序模型
// 具体参照tc-netem reorder 功能
type GapReorder struct {
	delay       time.Duration
	possibility float64
	correlation float64
	lastReorder bool
	gap         int
	count       int // 计数器，近似以gap为周期
	random      *rand.Rand
}

var _ Reorder = &GapReorder{}

func NewGapReorder(delay time.Duration, possibility,
	correlation float64, gap int, random *rand.Rand) node2.PacketHandler {
	reorder := &GapReorder{
		delay:       delay,
		possibility: possibility,
		correlation: correlation,
		lastReorder: false, // 默认首包之前的包没有reorder
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

func (gr *GapReorder) PacketHandler(*base.Packet, *base.PacketQueue) (time.Duration, bool) {
	return gr.Reorder(), false
}
