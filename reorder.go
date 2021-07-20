package networksimulator

import (
	"math/rand"
	"time"
)

// Reorder 接口，目前有不乱序，概率乱序，gap乱序三种模型
type Reorder interface {
	// 返回需要增加的延迟时间 (可能为负值)
	// 通过将该包更"早"发送实现乱序
	Reorder(rand *rand.Rand) float32
}

// NoneReorder 一律不乱序，不减少延迟时间
type NoneReorder struct {
	rand *rand.Rand
}

var _ Reorder = &NoneReorder{}

func NewNoneReorder(rand *rand.Rand) PacketHandler {
	reorder := &NoneReorder{rand: rand}
	return reorder.PacketHandler
}

func (nr *NoneReorder) Reorder(*rand.Rand) float32 {
	return 0
}

func (nr *NoneReorder) PacketHandler(*Packet, *PacketQueue) (delayTime time.Duration, isLoss bool) {
	return time.Microsecond * time.Duration(nr.Reorder(nr.rand)), false
}

// NormalReorder 概率乱序模型
// 以 correlation 的概率和上一个包乱序情况相同
// 剩余情况以 percent 概率乱序
type NormalReorder struct {
	avgDelay    float32
	percent     float32
	correlation float32
	lastReorder bool
	rand        *rand.Rand
}

var _ Reorder = &NormalReorder{}

func NewNormalReorder(avgDelay float32, percent float32,
	correlation float32, rand *rand.Rand) PacketHandler {
	reorder := &NormalReorder{
		avgDelay:    avgDelay,
		percent:     percent,
		correlation: correlation,
		lastReorder: false, // 默认首包之前的包没有reorder
		rand:        rand,
	}
	return reorder.PacketHandler
}

func (nr *NormalReorder) Reorder(rand *rand.Rand) float32 {
	if rand.Float32() < nr.correlation {
		if nr.lastReorder {
			return -nr.avgDelay
		} else {
			return 0
		}
	}
	if rand.Float32() >= nr.percent {
		nr.lastReorder = false
		return 0
	} else {
		nr.lastReorder = true
		return -nr.avgDelay
	}
}

func (nr *NormalReorder) PacketHandler(*Packet, *PacketQueue) (delayTime time.Duration, isLoss bool) {
	return time.Microsecond * time.Duration(nr.Reorder(nr.rand)), false
}

// GapReorder gap 乱序模型
// 具体参照tc-netem reorder 功能
type GapReorder struct {
	avgDelay    float32
	percent     float32
	correlation float32
	lastReorder bool
	gap         int
	count       int // 计数器，近似以gap为周期
	rand        *rand.Rand
}

var _ Reorder = &GapReorder{}

func NewGapReorder(avgDelay float32, percent float32,
	correlation float32, gap int, rand *rand.Rand) PacketHandler {
	reorder := &GapReorder{
		avgDelay:    avgDelay,
		percent:     percent,
		correlation: correlation,
		lastReorder: false, // 默认首包之前的包没有reorder
		gap:         gap,
		count:       0,
		rand:        rand,
	}
	return reorder.PacketHandler
}

func (gr *GapReorder) Reorder(rand *rand.Rand) float32 {
	if rand.Float32() < gr.correlation {
		if gr.lastReorder {
			return -gr.avgDelay
		} else {
			return 0
		}
	}

	gr.count++
	if gr.count < gr.gap {
		gr.lastReorder = false
		return 0
	}
	if rand.Float32() < gr.percent {
		gr.count = 0
		gr.lastReorder = true
		return -gr.avgDelay
	} else {
		gr.lastReorder = false
		return 0
	}
}

func (gr *GapReorder) PacketHandler(*Packet, *PacketQueue) (delayTime time.Duration, isLoss bool) {
	return time.Microsecond * time.Duration(gr.Reorder(gr.rand)), false
}
