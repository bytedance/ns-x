package main

import (
	"byte-ns"
	"byte-ns/base"
	"byte-ns/math"
	"byte-ns/node"
	"math/rand"
)

func main() {
	source := rand.NewSource(0)
	random := rand.New(source)
	helper := ns_x.NewBuilder()
	callback := func(packet *base.SimulatedPacket) {
		println("emit packet")
		println(packet.String())
	}
	n1 := node.NewChannelNode("entry1", 0, callback, math.NewRandomLoss(0.1, random))
	network, nodes := helper.
		Chain().
		Node(n1).
		Node(node.NewRestrictNode("", 0, nil, 1.0, 1024.0, 8192, 20)).
		Node(node.NewEndpointNode("endpoint")).
		Chain().
		Node(node.NewChannelNode("entry2", 0, callback, math.NewRandomLoss(0.1, random))).
		NodeByName("endpoint").
		Build(1, 10000, 10)
	network.Start()
	defer network.Stop()
	entry1 := nodes["entry1"]
	entry2 := nodes["entry2"]
	endpoint := nodes["endpoint"].(*node.EndpointNode)
	for i := 0; i < 20; i++ {
		entry1.Send([]byte{0x01, 0x02})
	}
	for i := 0; i < 20; i++ {
		entry2.Send([]byte{0x01, 0x02})
	}
	count := 0
	for {
		packet := endpoint.Receive()
		if packet != nil {
			count++
			println("receive packet")
			println(packet.String())
			println("total", count, "packets received")
		}
	}
}
