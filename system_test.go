package byte_ns

import (
	"math/rand"
	"testing"
)

func TestBasic(t *testing.T) {
	endpoint := NewEndpoint()
	source := rand.NewSource(0)
	random := rand.New(source)
	l := NewRandomLoss(0.32, random)
	node := NewChannel(0, func(packet *SimulatedPacket) {
		println("Emit packet ", packet.String())
	}, l)
	node.Next = endpoint
	nodes := []Node{endpoint, node}
	network := NewNetwork(nodes)
	network.Start()
	defer network.Stop()
	node.Send(&Packet{Data: []byte{0x01, 0x02}})
	node.Send(&Packet{Data: []byte{0x02, 0x03}})
	node.Send(&Packet{Data: []byte{0x03, 0x04}})
	for {
		packet := endpoint.Receive()
		if packet != nil {
			println("receive packet ", packet.String())
		}
	}
}
