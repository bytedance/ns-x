# ByteNS

An easy-to-use, flexible library to simulate network behavior, written mainly in Go.

## Feature

* Build highly-customizable and scalable network topology upon basic nodes.
* Simulate loss, delay, etc. on any nodes by any given parameters and models.
* Collect data from each and every node in detail.
* Cross-platform.

## Introduction

#### Concept

* Network: a topological graph consist of *node*s, reflecting a real-world network for packets to transfer through.
* Node: a physical or logical device in the *network* allowing *packet*s to enter and leave. A *node* usually chains other *Node*s to build the network.
* Packet: simulated data packets transferring between *node*s, carrying the actual packet data with additional simulator information.
* Send: (*packet*s) to enter a node, waiting to *emit*.
* Emit: (*packet*s) to leave a *node* toward the next chained *node*.

#### Prerequisites

- `Go mod` must be supported and enabled.  
- A platform-specific `binary/*/libtime.a` library is required by cgo for high resolution timer. Its Windows, Linux, and Darwin binaries are pre-built. Compile the library manually if running on another arch/os. (See <a href = "#compile">compile</a> section)    

#### Usage

Follow three steps: building network, simulating, and collecting data.

##### Building network

The network is built by nodes and chains (i.e., edges in the network graph). Normally a chain connects only two nodes, each on one end. For special cases a chain may have multiple node or no node on its ends.  

Nodes are highly customizable, and some typical nodes are pre-defined:   

* Broadcast: a node transfers packet from one source to multiple targets.
* Channel: a node delays, losses or reorders packets passing by.
* Endpoint: a node only accepts incoming packets, usually acting as the end of a chain.
* Gather: a node gathers packets from multiple sources to a single target.
* Restrict: a node limits pps or bps by dropping packets when its internal buffer overflows.
* Scatter: a node selects which node the incoming packet should be route to according to a given rule.

##### Simulating

Once the network is built, packets can be sent into any entry nodes and received from any exit nodes.

##### Collecting Data

Data could be collected by callback function `node.OnEmitCallback()`. Any further analyses to the collected data could be done after the simulation. Also note that time-costing callbacks would slow down the simulation so keep an eye on its performance.

#### Example

A network builder is also provided in order to describe the whole network conveniently.

Following is an example of sending packets through a simulated channel with `32%` packet loss.

```go
package main

import (
	"byte_ns"
	"math/rand"
)

func main() {
	source := rand.NewSource(0)
	random := rand.New(source)
	helper := byte_ns.NewBuilder()
	callback := func(packet *byte_ns.SimulatedPacket) {
		println("emit packet", packet)
	}
	n1 := byte_ns.NewChannel("entry1", 0, callback, byte_ns.NewRandomLoss(0.3, random))
	network, nodes := helper.
		Chain().
		Node(n1).
		Node(byte_ns.NewRestrict("", 0, nil, 1.0, 1024.0, 4096, 5)).
		Node(byte_ns.NewEndpoint("endpoint")).
		Chain().
		Node(byte_ns.NewChannel("entry2", 0, callback, byte_ns.NewRandomLoss(0.1, random))).
		NodeByName("endpoint").
		Build()
	network.Start()
	defer network.Stop()
	entry1 := nodes["entry1"]
	entry2 := nodes["entry2"]
	endpoint := nodes["endpoint"].(*byte_ns.Endpoint)
	entry1.Send(&byte_ns.Packet{Data: []byte{0x01, 0x02}})
	entry1.Send(&byte_ns.Packet{Data: []byte{0x02, 0x03}})
	entry1.Send(&byte_ns.Packet{Data: []byte{0x03, 0x04}})
	entry2.Send(&byte_ns.Packet{Data: []byte{0x03, 0x04}})
	entry2.Send(&byte_ns.Packet{Data: []byte{0x03, 0x04}})
	for {
		packet := endpoint.Receive()
		if packet != nil {
			println("receive packet ", packet.String())
		}
	}
}
```

#### Compile libtime<span id="compile"/>

The following library is built successfully on Go v1.16.5, cmake v3.21.0, clang v12.0.5, with C++ 11.

```bash
cd cpp
cmake CMakeLists.txt
make
```

which generates file `libtime.a` under `cpp` directory.

To make the compiled library work, a tag *time_compiled* need to be added to go build.

```bash
go build -tags time_compiled
```

There is also a configuration file `cross-compile.cmake` for cross compiling the high resolution time library with little modification.

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
* implement commonly used protocol stack as a new node type
