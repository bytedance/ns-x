package main

import (
	"github.com/bytedance/ns-x/v2"
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/math"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/bytedance/ns-x/v2/tick"
	"math/rand"
	"time"
)

func echo() {
	now := time.Now()
	helper := ns_x.NewBuilder()
	random := rand.New(rand.NewSource(0))
	network, nodes := helper.
		Chain().
		NodeWithName("restrict 1", node.NewRestrictNode(node.WithBPSLimit(1024*1024, 4*1024*1024))).
		NodeWithName("channel 1", node.NewChannelNode(node.WithLoss(math.NewRandomLoss(0.1, random)))).
		Chain().
		NodeWithName("restrict 2", node.NewRestrictNode(node.WithPPSLimit(10, 50))).
		NodeWithName("channel 2", node.NewChannelNode(node.WithLoss(math.NewRandomLoss(0.3, random)))).
		Chain().
		NodeWithName("endpoint 1", node.NewEndpointNode()).
		Group("restrict 1", "channel 1").
		NodeWithName("endpoint 2", node.NewEndpointNode()).
		Chain().
		NodeOfName("endpoint 2").
		Group("restrict 2", "channel 2").
		NodeOfName("endpoint 1").
		Summary().
		Build()
	endpoint1 := nodes["endpoint 1"].(*node.EndpointNode)
	endpoint2 := nodes["endpoint 2"].(*node.EndpointNode)
	endpoint1.Receive(func(packet base.Packet, now time.Time) []base.Event {
		return base.Aggregate(endpoint1.Send(packet, now))
	})
	endpoint2.Receive(func(packet base.Packet, now time.Time) []base.Event {
		return base.Aggregate(endpoint2.Send(packet, now))
	})
	network.Run([]base.Event{endpoint1.Send(base.RawPacket("hello world"), now)}, tick.NewStepClock(now, time.Second), 30*time.Second)
	defer network.Wait()
	time.Sleep(time.Second)
}
