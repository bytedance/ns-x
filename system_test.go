package byte_ns

import (
	"math/rand"
	"runtime"
	"testing"
)

func TestBasic(t *testing.T) {
	source := rand.NewSource(0)
	random := rand.New(source)
	helper := NewBuilder()
	callback := func(packet *SimulatedPacket) {
		println("emit packet")
		println(packet.String())
	}
	n1 := NewChannel("entry1", 0, callback, NewRandomLoss(0.3, random))
	network, nodes := helper.
		Chain().
		Node(n1).
		Node(NewRestrict("", 0, nil, 1.0, 1024.0, 4096, 5)).
		Node(NewEndpoint("endpoint")).
		Chain().
		Node(NewChannel("entry2", 0, callback, NewRandomLoss(0.1, random))).
		NodeByName("endpoint").
		Build(runtime.NumCPU()/2, 10, 10)
	network.Start()
	defer network.Stop()
	entry1 := nodes["entry1"]
	entry2 := nodes["entry2"]
	endpoint := nodes["endpoint"].(*Endpoint)
	for i := 0; i < 20; i++ {
		entry1.Send(&Packet{Data: []byte{0x01, 0x02}})
	}
	for i := 0; i < 20; i++ {
		entry2.Send(&Packet{Data: []byte{0x01, 0x02}})
	}
	for {
		packet := endpoint.Receive()
		if packet != nil {
			println("receive packet")
			println(packet.String())
		}
	}
}
