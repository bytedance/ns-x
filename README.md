# Network Simulator

Network Simulator is designed as an easy-to-use, flexible library to simulate internet, written mainly in go.

### Feature

* Flexible to construct any network graph with customizable node at any scale
* Ability to collect any data from any node in the network graph
* Cross platform, can be used on any platform / architecture supports go and C++

### Introduction

#### Network

As the project name, network simulator simulates the real network. Real network are abstracted into simple math models and connections of them, these math models are called nodes in network simulator. A network in network simulator contains all the nodes, and transfer packets from one node to another node.

Each node has two methods: Send and Emit, Send means the packet enter the node while Emit means the packet leave the node and transfer to the next node if necessary. Each node has a packet buffer as well, all packets sent will be put in the buffer, the network will clear each node's buffer frequently, and emit them later. The buffer itself is implemented thread-safe and lock-free for high performance.

#### Simulated Packets

Simulated packets contains where is the packet (current node),  when the packet enters current node and plan to leave (sent time and emit time), whether the packet is lost, as well as the actual packet.

#### Node

Some widely used math models are already pre-defined:

* Broadcast: A broadcast node can transfer the same packet to multiple target.
* Channel: A channel node can delay, loss or reorder packets through it.
* Endpoint: An endpoint node is where packets can be received, this node is usually the end of a node chain.
* Gather: A gather node gather packets from multiple sources, and transfer them to a target.
* Restrict: A restrict node block following packets when reach the restrict, and drop following packets once the internal buffer overflow.
* Scatter: A scatter node transfer packets of a source to one of its targets selected by a user-defined rule.

### Usage

#### Install

The installation only requires to add it into go.mod.

The project use cgo to implement high resolution time, by default, binary for windows, linux and bsd with amd64 are prebuilt, other platforms or architectures need to compile it by self. (See the <a href = "#compile">compile</a> section)

#### Example

Following is an example of send packets through a simulated channel, with packet loss posibility of $32\%$​.

```go
package main

import (
	"math/rand"
	"network-simulator"
)

func main() {
	endpoint := networksimulator.NewEndpoint()
	source := rand.NewSource(0)
	random := rand.New(source)
	l := networksimulator.NewRandomLoss(0.32, random)
	n := networksimulator.NewChannel(endpoint, 0, func(packet *networksimulator.SimulatedPacket) {
		println("Emit packet ", packet.String())
	}, l)
	nodes := []networksimulator.Node{endpoint, n}
	network := networksimulator.NewNetwork(nodes)
	network.Start()
	defer network.Stop()
	n.Send(&networksimulator.Packet{Data: []byte{0x01, 0x02}})
	n.Send(&networksimulator.Packet{Data: []byte{0x02, 0x03}})
	n.Send(&networksimulator.Packet{Data: []byte{0x03, 0x04}})
	for {
		packet := endpoint.Receive()
		if packet != nil {
			println("receive packet ", packet.String())
		}
	}
}
```

#### <div id="compile">Compile

The build environment tested is go v1.16.5, with cmake v3.21.0, clang v12.0.5 and C++ 11.

The project use cgo to implement high resolution time, to compile it, goto core/cpp and use cmake to build.

```bash
cd core/cpp
cmake CMakeLists.txt
make
```

This should generate a file named $libtime.a$​ under the cpp directory.

To make the compiled library work, a tag *time_compiled* need to be added to go build.

### Design

#### Main Loop

The main loop maintains a thread-local heap, in order to sort the packets.

The loop is separated into two parts: fetch and drain.

* fetch: The main loop clear the packet buffer of all the nodes, and put these packets into the heap.
* drain: The main loop drain the heap until the heap only contains packets with emit time after the current time.

By now, the main loop lock a single os thread, but in the future, the main loop may run on a fork join pool.

#### High Resolution Time

Since the main loop need to access current time with high resolution as well as low cost, the standard time library of go is not enough. (internal system call, accurate but with high cost, update not timely)

Currently, high resolution time is a wrapper of C++ time library. The core design is use system time and steady time together. The system time means time retrieved through system call, while the steady time is usually a counter of CPU cycles. The system time is accurate but with lower resolution and higher cost, the steady time is not so accurate (due to turbo of CPU) but with highest resolution in theory. Once trying to fetch the time, it's checked whether enough time has passed by since the last alignment. If so, an align will be performed immediately. The align operation itself is thread safe by a lock, but another double check will guarantee low cost.

