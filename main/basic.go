package main

// Example of how to use the simulator basically

import (
	"github.com/bytedance/ns-x/v2"
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/math"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/bytedance/ns-x/v2/tick"
	"go.uber.org/atomic"
	"math/rand"
	"time"
)

func basic() {
	source := rand.NewSource(0)
	random := rand.New(source)
	helper := ns_x.NewBuilder()
	callback := func(packet base.Packet, source, target base.Node, now time.Time) {
		println("emit packet")
	}
	n1 := node.NewEndpointNode()
	t := time.Now()
	network, nodes := helper.
		Chain().
		NodeWithName("entry1", n1).
		Node(node.NewChannelNode(node.WithTransferCallback(callback), node.WithLoss(math.NewRandomLoss(0.1, random)))).
		Node(node.NewRestrictNode(node.WithPPSLimit(1, 20))).
		NodeWithName("endpoint", node.NewEndpointNode()).
		Chain().
		NodeWithName("entry2", node.NewEndpointNode()).
		Node(node.NewChannelNode(node.WithTransferCallback(callback), node.WithLoss(math.NewRandomLoss(0.1, random)))).
		NodeOfName("endpoint").
		Summary().
		Build()
	entry1 := nodes["entry1"].(*node.EndpointNode)
	entry2 := nodes["entry2"].(*node.EndpointNode)
	endpoint := nodes["endpoint"].(*node.EndpointNode)
	count := atomic.NewInt64(0)
	endpoint.Receive(func(packet base.Packet, now time.Time) []base.Event {
		if packet != nil {
			count.Inc()
			println("receive packet at", now.String())
			println("total", count.Load(), "packets received")
		}
		return nil
	})
	total := 20
	events := make([]base.Event, 0, total*2)
	for i := 0; i < 20; i++ {
		events = append(events, entry1.Send(base.RawPacket([]byte{0x01, 0x02}), t))
	}
	for i := 0; i < 20; i++ {
		events = append(events, entry2.Send(base.RawPacket([]byte{0x01, 0x02}), t))
	}
	event := base.NewPeriodicEvent(func(now time.Time) []base.Event {
		for i := 0; i < 10; i++ {
			_ = rand.Int()
		}
		return nil
	}, time.Second, t)
	events = append(events, event)
	network.Run(events, tick.NewStepClock(t, time.Second), 300*time.Second)
	defer network.Wait()
}
