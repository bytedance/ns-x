package node

import (
	"math/rand"
	"network-simulator/core"
	"time"
)

// Delay 接口
// 目前有无分布，正态分布，均匀分布三种实现
type Delay interface {
	Delay() time.Duration   // 返回具体延迟（包含jitter）
	Average() time.Duration // 返回参数 average
	Jitter() time.Duration  // 返回 jitter 结果
}

// basicDelay 是 Delay 共同约定使用的变量
type basicDelay struct {
	average time.Duration // 平均延迟 单位 ms
	jitter  time.Duration // 单位 ms
	random  *rand.Rand    // 节点统一的随机数生成器
}

// Average 返回参数 average
func (bd *basicDelay) Average() time.Duration {
	return bd.average
}

// Jitter 返回参数 jitter
func (bd *basicDelay) Jitter() time.Duration {
	jitter := float64(bd.jitter.Microseconds())
	return time.Duration((2*bd.random.Float64()-1)*jitter) * time.Microsecond
}

// FixedDelay 无分布，每次均为 average
type FixedDelay struct {
	*basicDelay
}

var _ Delay = &FixedDelay{}

// NewFixedDelay 创建一个无分布延迟模型处理函数，average，jitter 单位是 ms
func NewFixedDelay(average time.Duration, jitter time.Duration) PacketHandler {
	delay := &FixedDelay{
		&basicDelay{
			average: average,
			jitter:  jitter,
		},
	}
	return delay.PacketHandler
}

func (nd *FixedDelay) Delay() time.Duration {
	return nd.Average() + nd.Jitter()
}

func (nd *FixedDelay) PacketHandler(*core.Packet, *core.PacketQueue) (time.Duration, bool) {
	return nd.Delay(), false
}

// NormalDelay 正态分布，大体在 [average-3*sigma,average+3*sigma] 范围
type NormalDelay struct {
	*basicDelay
	sigma time.Duration
}

var _ Delay = &NormalDelay{}

func NewNormalDelay(average, jitter, sigma time.Duration, random *rand.Rand) PacketHandler {
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

func (nd *NormalDelay) PacketHandler(*core.Packet, *core.PacketQueue) (time.Duration, bool) {
	return nd.Delay(), false
}

// UniformDelay 均匀分布，[0,2*average) 范围内
type UniformDelay struct {
	*basicDelay
}

var _ Delay = &UniformDelay{}

func NewUniformDelay(average, jitter time.Duration, random *rand.Rand) PacketHandler {
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

func (ud *UniformDelay) PacketHandler(*core.Packet, *core.PacketQueue) (time.Duration, bool) {
	return ud.Delay(), false
}
