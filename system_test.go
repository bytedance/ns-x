package networksimulator

import (
	"math/rand"
	"network-simulator/core"
	"network-simulator/node"
	"testing"
)

func TestBasic(t *testing.T) {
	endpoint := node.NewEndpoint()
	source := rand.NewSource(0)
	random := rand.New(source)
	l := node.NewRandomLoss(0.32, random)
	n := node.NewChannel(endpoint, 0, func(packet *core.SimulatedPacket) {
		println("Emit packet ", packet.String())
	}, l)
	nodes := []core.Node{endpoint, n}
	network := core.NewNetwork(nodes)
	network.Start()
	defer network.Stop()
	n.Send(&core.Packet{Data: []byte{0x01, 0x02}})
	n.Send(&core.Packet{Data: []byte{0x02, 0x03}})
	n.Send(&core.Packet{Data: []byte{0x03, 0x04}})
	for {
		packet := endpoint.Receive()
		if packet != nil {
			println("receive packet ", packet.String())
		}
	}
}
