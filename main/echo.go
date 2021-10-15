package main

// Example of a duplex network, where to endpoints echo to each other

import (
	"github.com/bytedance/ns-x/v2"
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/math"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/bytedance/ns-x/v2/tick"
	"time"
)

func main() {
	now := time.Now()
	helper := ns_x.NewBuilder()
	network, nodes := helper.
		Chain().
		NodeWithName("restrict 1", node.NewRestrictNode(node.WithBPSLimit(1024*1024, 4*1024*1024))).
		NodeWithName("channel 1", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(150*time.Millisecond)))).
		Chain().
		NodeWithName("restrict 2", node.NewRestrictNode(node.WithPPSLimit(10, 50))).
		NodeWithName("channel 2", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(200*time.Millisecond)))).
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
		println("endpoint 1 receive:", string(packet.(base.RawPacket)), "at", now.String())
		return base.Aggregate(endpoint1.Send(packet, now))
	})
	endpoint2.Receive(func(packet base.Packet, now time.Time) []base.Event {
		println("endpoint 2 receive:", string(packet.(base.RawPacket)), "at", now.String())
		return base.Aggregate(endpoint2.Send(packet, now))
	})
	network.Run([]base.Event{endpoint1.Send(base.RawPacket("hello world"), now)}, tick.NewStepClock(now, time.Second), 30*time.Second)
	defer network.Wait()
}
