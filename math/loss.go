package math

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"math/rand"
	"time"
)

// Loss in the network
type Loss interface {
	Loss() bool // the packet Loss if true
}

// NoneLoss no loss
type NoneLoss struct {
}

var _ Loss = &NoneLoss{}

func NewNoneLoss() node.PacketHandler {
	loss := &NoneLoss{}
	return loss.PacketHandler
}

func (nl *NoneLoss) Loss() bool {
	return false
}

func (nl *NoneLoss) PacketHandler(base.Packet) (time.Duration, bool) {
	return 0, nl.Loss()
}

// RandomLoss loss with the given possibility
type RandomLoss struct {
	possibility float64
	random      *rand.Rand
}

var _ Loss = &RandomLoss{}

func NewRandomLoss(possibility float64, random *rand.Rand) node.PacketHandler {
	loss := &RandomLoss{
		possibility: possibility,
		random:      random,
	}
	return loss.PacketHandler
}

func (rl *RandomLoss) Loss() bool {
	return rl.random.Float64() < rl.possibility
}

func (rl *RandomLoss) PacketHandler(base.Packet) (time.Duration, bool) {
	return 0, rl.Loss()
}

// GilbertLoss Gilbert Loss Model
type GilbertLoss struct {
	s1Loss, s1Transit, s2Loss, s2Transit float64
	gilbertState                         int
	random                               *rand.Rand
}

var _ Loss = &GilbertLoss{}

func NewGilbertLoss(s1Loss, s1Transit, s2Loss, s2Transit float64) node.PacketHandler {
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

func (gl *GilbertLoss) PacketHandler(base.Packet) (time.Duration, bool) {
	return 0, gl.Loss()
}
