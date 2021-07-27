package byte_ns

import (
	"math/rand"
	"time"
)

// Loss 接口，
// 目前有不丢包，单一概率丢包模型和 Gilbert 模型丢包 三种具体实现
type Loss interface {
	Loss() bool // 若返回 true 则丢包
}

// NoneLoss 不丢包
type NoneLoss struct {
}

var _ Loss = &NoneLoss{}

func NewNoneLoss() PacketHandler {
	loss := &NoneLoss{}
	return loss.PacketHandler
}

func (nl *NoneLoss) Loss() bool {
	return false
}

func (nl *NoneLoss) PacketHandler(*Packet, *PacketQueue) (time.Duration, bool) {
	return 0, nl.Loss()
}

// RandomLoss 单一概率丢包模型
type RandomLoss struct {
	possibility float64
	random      *rand.Rand
}

var _ Loss = &RandomLoss{}

func NewRandomLoss(possibility float64, random *rand.Rand) PacketHandler {
	loss := &RandomLoss{
		possibility: possibility,
		random:      random,
	}
	return loss.PacketHandler
}

func (rl *RandomLoss) Loss() bool {
	return rl.random.Float64() < rl.possibility
}

func (rl *RandomLoss) PacketHandler(*Packet, *PacketQueue) (time.Duration, bool) {
	return 0, rl.Loss()
}

// GilbertLoss Gilbert丢包模型
// 有两个状态，分别有一个丢包概率和一个状态变迁概率
type GilbertLoss struct {
	s1Loss, s1Transit, s2Loss, s2Transit float64
	gilbertState                         int
	random                               *rand.Rand
}

var _ Loss = &GilbertLoss{}

func NewGilbertLoss(s1Loss, s1Transit, s2Loss, s2Transit float64) PacketHandler {
	loss := &GilbertLoss{
		s1Loss:       s1Loss,
		s1Transit:    s1Transit,
		s2Loss:       s2Loss,
		s2Transit:    s2Transit,
		gilbertState: 0,
	}
	return loss.PacketHandler
}

func (gl *GilbertLoss) Loss() bool {
	if gl.gilbertState == 0 {
		if gl.random.Float64() < gl.s1Transit {
			gl.gilbertState = 1
		}
		return gl.random.Float64() < gl.s1Loss
	}
	if gl.random.Float64() < gl.s2Transit {
		gl.gilbertState = 0
	}
	return gl.random.Float64() < gl.s2Loss
}

func (gl *GilbertLoss) PacketHandler(*Packet, *PacketQueue) (time.Duration, bool) {
	return 0, gl.Loss()
}
