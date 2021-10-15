package main

// Example of how to customize route rule

import (
	"github.com/bytedance/ns-x/v2"
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/math"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/bytedance/ns-x/v2/tick"
	"time"
)

func route() {
	helper := ns_x.NewBuilder()
	t := time.Now()
	routeTable := make(map[base.Node]base.Node)
	ipTable := make(map[string]base.Node)
	scatter := node.NewScatterNode(node.WithRouteSelector(func(packet base.Packet, nodes []base.Node) base.Node {
		if p, ok := packet.(*packetWithNode); ok {
			return routeTable[p.destination]
		}
		panic("no route to host")
	}))
	client := node.NewEndpointNode()
	network, nodes := helper.
		Chain().
		Node(client).
		Node(scatter).
		NodeWithName("route1", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(time.Millisecond*200)))).
		NodeWithName("server1", node.NewEndpointNode()).
		Chain().
		Node(client).
		Node(scatter).
		NodeWithName("route2", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(time.Millisecond*300)))).
		NodeWithName("server2", node.NewEndpointNode()).
		Build()
	server1 := nodes["server1"].(*node.EndpointNode)
	server2 := nodes["server2"].(*node.EndpointNode)
	route1 := nodes["route1"]
	route2 := nodes["route2"]
	routeTable[server1] = route1
	routeTable[server2] = route2
	ipTable["192.168.0.1"] = server1
	ipTable["192.168.0.2"] = server2
	server1.Receive(react1) // server 1 should receive after 1-second send delay + 200 milliseconds channel delay
	server2.Receive(react2) // server 2 should receive after 2-second send delay + 200 milliseconds channel delay
	sender := createSender(client, ipTable)
	events := make([]base.Event, 0)
	events = append(events, sender(base.RawPacket([]byte{}), "192.168.0.1", t.Add(time.Second*1))) // send to server1 after 1 second
	events = append(events, sender(base.RawPacket([]byte{}), "192.168.0.2", t.Add(time.Second*2))) // send to server2 after 2 second
	network.Run(events, tick.NewStepClock(t, time.Millisecond), 300*time.Second)
	defer network.Wait()
}

func react1(packet base.Packet, now time.Time) []base.Event {
	println("server 1 receive at", now.String())
	return nil
}

func react2(packet base.Packet, now time.Time) []base.Event {
	println("server 2 receive at", now.String())
	return nil
}

type packetWithNode struct {
	base.Packet
	source, destination base.Node
}

type sender func(packet base.Packet, ip string, t time.Time) base.Event

func createSender(client *node.EndpointNode, ipTable map[string]base.Node) sender {
	return func(packet base.Packet, ip string, t time.Time) base.Event {
		return client.Send(&packetWithNode{packet, client, ipTable[ip]}, t)
	}
}
