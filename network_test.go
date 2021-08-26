package ns_x

import (
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/tick"
	"math/rand"
	"testing"
	"time"
)

func nop(time.Time) []base.Event {
	return nil
}

func BenchmarkEventLoop(b *testing.B) {
	network := NewNetwork([]base.Node{}, tick.NewRealClock())
	now := time.Now()
	queue := base.NewEventQueue(time.Second)
	for i := 0; i < b.N; i++ {
		queue.Enqueue(base.NewFixedEvent(nop, now.Add(-time.Duration(rand.Int()%20)*time.Second)))
	}
	b.ResetTimer()
	network.eventLoop(queue)
}
