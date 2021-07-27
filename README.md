# Network Simulator

Network Simulator is designed as an easy-to-use, flexible library to simulate internet, written mainly in go.

## Feature

* Flexible to construct any network graph with customizable node at any scale
* Ability to collect any data from any node in the network graph
* Cross-platform, can be used on any platform / architecture supports go and C++

## Introduction

#### Concept

* Network: Simulated network where packets transfer through
* Node: Simulated physical or logical device in the network
* Packets: Simulated packet, contains the actual packet with some more information
* Send: A packet sent to a node means the packet enters the node, waiting to emit
* Emit: A packet is emitted means the packet leaves the node, transfer to the next node

#### Install

The installation only requires to add it into go.mod.

The project use cgo to implement high resolution time, by default, binary for windows, linux and bsd with amd64 are pre-built, other platforms or architectures need to compile it by self. (See the <a href = "#compile">compile</a> section)

#### Usage

The simulation can be separated into three parts: build model, simulate and collect data.

##### Build Model

Node can be connected to other nodes and form the network, normal node can only have exactly one node as it's next node, but some kind of nodes can have multiple next nodes, or even no next nodes.

Each node has two methods: Send and Emit, Send means the packet enter the node while Emit means the packet leave the node and transfer to the next node if necessary. 

Each node has a callback called when a packet emitted, where users can collect data.

Node can be customized with high flexibility, but some widely used nodes are already pre-defined:

* Broadcast: A broadcast node can transfer the same packet to multiple target.
* Channel: A channel node can delay, loss or reorder packets through it.
* Endpoint: An endpoint node is where packets can be received, this node is usually the end of a node chain.
* Gather: A gather node gather packets from multiple sources, and transfer them to a target.
* Restrict: A restrict node block following packets when reach the restrict, and drop following packets once the internal buffer overflow.
* Scatter: A scatter node transfer packets of a source to one of its targets selected by a user-defined rule.

##### Simulate

Once the network is built, packets can be sent to the entry and received at the exit. The simulated network and nodes will act as defined to transfer data, until simulation finished.

##### Collect Data

Data are collected during the simulation through registered callback. However, only necessary work can be done in the callback, or the simulation would be affected. Once simulation finished, further analyze can be done on the collected data.

#### Example

Following is an example of sending packets through a simulated channel, with packet-loss possibility of $$32\%$$.

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

#### Compile<span id="compile"/>

The build environment tested is go v1.16.5, with cmake v3.21.0, clang v12.0.5 and C++ 11.

The project use cgo to implement high resolution time under the cpp directory, to compile it, go to the directory and use cmake to build.

```bash
cd cpp
cmake CMakeLists.txt
make
```

This should generate a file named $libtime.a$ under the cpp directory.

To make the compiled library work, a tag *time_compiled* need to be added to go build.

There is also a configuration named *cross-compile.cmake*, which can be used to cross compile the high resolution time library with little modification.

## Design

#### Architecture

Each node has a packet buffer, once a packet is sent to the node, it will be put in the buffer. The buffer itself is implemented thread-safe and lock-free for high performance.

There is a global main loop host by the network, which clears each node's buffer and decide when to emit packets.

#### Main Loop

The main loop maintains a thread-local heap, in order to sort the packets.

The loop is separated into two parts: fetch and drain.

* fetch: The main loop clear the packet buffer of all the nodes, and put these packets into the heap.
* drain: The main loop drain the heap until the heap only contains packets with emit time after the current time.

By now, the main loop lock a single os thread, but in the future, the main loop may run on a fork join pool.

#### High Resolution Time

Time is of vital significance in simulations, it directly decides the accuracy of simulation.

Since the simulation needs to access current time with high resolution and low cost, the standard time library of go is not enough. (internal system call, accurate but with high cost, update not timely)

Currently, high resolution time is a wrapper of C++ time library. The core design is use system time and steady time together. The system time means time retrieved through system call, while the steady time is usually a counter of CPU cycles. The system time is accurate but with lower resolution and higher cost, the steady time is not so accurate (due to turbo of CPU) but with the highest resolution in theory. Once trying to fetch the time, it's checked whether enough time has passed by since the last alignment. If so, an alignment will be performed immediately. The align operation itself is thread safe by a lock, but another double check will guarantee low cost.

## Contribution

#### Future work

* parallelize main loop
* implement commonly used protocol stack as node

