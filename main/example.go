package main

import (
	"go.uber.org/atomic"
	"math/rand"
	"ns-x"
	"ns-x/base"
	"ns-x/math"
	"ns-x/node"
	"time"
)

func main() {
	source := rand.NewSource(0)
	random := rand.New(source)
	helper := ns_x.NewBuilder()
	callback := func(packet base.Packet, target base.Node, now time.Time) {
		println("emit packet")
	}
	n1 := node.NewEndpointNode("entry1", nil)
	network, nodes := helper.
		Chain().
		Node(n1).
		Node(node.NewChannelNode("", callback, math.NewRandomLoss(0.1, random))).
		Node(node.NewRestrictNode("", nil, 1.0, 1024.0, 8192, 20)).
		Node(node.NewEndpointNode("endpoint", nil)).
		Chain().
		Node(node.NewEndpointNode("entry2", nil)).
		Node(node.NewChannelNode("", callback, math.NewRandomLoss(0.1, random))).
		NodeByName("endpoint").
		Build(1, 10000, 10)
	entry1 := nodes["entry1"].(*node.EndpointNode)
	entry2 := nodes["entry2"].(*node.EndpointNode)
	endpoint := nodes["endpoint"].(*node.EndpointNode)
	count := atomic.NewInt64(0)
	endpoint.Receive(func(packet base.Packet, now time.Time) {
		if packet != nil {
			count.Inc()
			println("receive packet at", now.String())
			println("total", count.Load(), "packets received")
		}
	})
	for i := 0; i < 20; i++ {
		network.Event(entry1.Send(base.RawPacket([]byte{0x01, 0x02})))
	}
	for i := 0; i < 20; i++ {
		network.Event(entry2.Send(base.RawPacket([]byte{0x01, 0x02})))
	}
	network.Start()
	defer network.Stop()
	time.Sleep(30 * time.Second)
}
