# ns-x

[![Go](https://github.com/bytedance/ns-x/actions/workflows/go.yml/badge.svg)](https://github.com/bytedance/ns-x/actions/workflows/go.yml)
[![CodeQL](https://github.com/bytedance/ns-x/actions/workflows/codeql.yml/badge.svg)](https://github.com/bytedance/ns-x/actions/workflows/codeql.yml)

An easy-to-use, flexible **network simulator** library, written mainly in Go.

## Feature

* Build highly-customizable and scalable network topology upon basic nodes.
* Simulate loss, delay, etc. on any nodes by any given parameters and models.
* Well-defined simulation model with definite math properties.
* Collect data from each and every node in detail.
* Cross-platform.

## Introduction

#### Concept

* Network: a topological graph consist of *node*s, reflecting a real-world network for *packets* to transfer through.
* Node: a physical or logical device in the *network* deciding what to do when a packet pass by. A *node* usually *connects* to other *node*s.
* Event: some action to be done at a given time point.
* Packet: simulated data packets transferring between *node*s.
* Transfer: the behavior of *node* when a *packet* pass by.

#### Prerequisites

- `go mod` must be supported and enabled.

#### Usage

Follow three steps: building network, starting network simulation, and collecting data.

##### Building network

The network is built by *node*s and *edge*s. Normally an *edge* connects only two nodes, each on one end. For special cases a chain may have multiple node or no node on its ends.

While nodes are highly customizable, some typical nodes are pre-defined as follows:

<!-- note that mermaid compilation of GitHub action only supports code blocks with no indents -->

* Broadcast: a node transfers packet from one source to multiple targets.

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 1~](/.resources/README-md-1.svg)
<details>
  <summary>Mermaid markup</summary>

```mermaid
graph LR;
    In --> Broadcast --> Out1;
    Broadcast --> Out2;
    Broadcast --> Out3;
```

</details>
<!-- generated by mermaid compile action - END -->

* Channel: a node delays, losses or reorders packets passing by.

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 2~](/.resources/README-md-2.svg)
<details>
  <summary>Mermaid markup</summary>

```mermaid
graph LR;
    In --> Channel -->|Loss & Delay & Reorder| Out;
```

</details>
<!-- generated by mermaid compile action - END -->

* Endpoint: a node where to send and receive packets, usually acting as the endpoint of a chain.

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 3~](/.resources/README-md-3.svg)
<details>
  <summary>Mermaid markup</summary>

```mermaid
graph LR;
   Nodes... --> Endpoint;
```

</details>
<!-- generated by mermaid compile action - END -->

* Gather: a node gathers packets from multiple sources to a single target.

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 4~](/.resources/README-md-4.svg)
<details>
  <summary>Mermaid markup</summary>

```mermaid
graph LR;
    In1 --> Gather ==> Out;
    In2 --> Gather;
    In3 --> Gather;
```

</details>
<!-- generated by mermaid compile action - END -->

* Restrict: a node limits pps or bps by dropping packets when its internal buffer overflows.

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 5~](/.resources/README-md-5.svg)
<details>
  <summary>Mermaid markup</summary>

```mermaid
graph LR;
    In --> Restrict -->|Restricted Speed| Out;
```

</details>
<!-- generated by mermaid compile action - END -->

* Scatter: a node selects which node the incoming packet should be route to according to a given rule.

<!-- generated by mermaid compile action - START -->
![~mermaid diagram 6~](/.resources/README-md-6.svg)
<details>
  <summary>Mermaid markup</summary>

```mermaid
graph LR;
    In --> Scatter -.-> Out1;
    Scatter -->|Selected Route| Out2;
    Scatter -.-> Out3;
```

</details>
<!-- generated by mermaid compile action - END -->

After all necessary *node*s created, *connect* them to build the network. To do so, just set the *next node* correctly for each *node* to declare the *edge*.

ns-x also provides a *builder* to facilitate the process. Instead of connecting *edge*s, it builds the network by connecting all *path*s in one line of code. *Path*, aka *chain*, is similar to the *path* concept in graph theory, representing a route along the *edge*s of a graph.

* `Chain()`: saves current chain (*path*) and in order to describe another chain.
* `Node()`: appends a given node to current chain.
* `NodeWithName()`: same as *Node*, with a customizable name to refer to later.
* `NodeByName()`: finds (refer to) a node with the given name, and appends it to current chain.
* `NodeGroup()`: given some nodes, perform *Node* operation on each of them in order.
* `NodeGroupWithName()`: same as *NodeGroup*, with a customizable name.
* `NodeGroupByName()`: finds a group with the given name, then perform *NodeGroup* operation on it.
* `Build()`: actually connect the previously described *chain*s to finally build the network. Note that connections of nodes outside the builder will be overwritten

##### Starting Network Simulation

Once the network is built, start running it so packets can be sent into and received from any *endpoint* nodes.

##### Collecting Data

Data could be collected by callback function `node.OnTransferCallback()`. Also note that time-costing callbacks would slow down the simulation and lead to error of result, so it is highly recommended only collecting data in the callbacks.Further analyses should be done after the simulation.

#### Property

Some properties are guaranteed:

* Order: if any event e at time point t, only generate events at time point not before t, then the handling order of two events at different time point is guaranteed, and the order of events at same time point is undetermined.
* Accurate: each event will be handled at the given time point exactly in simulate clock, and the difference between simulate clock and real clock is as small as possible, usually some microseconds.

Some other properties of each kind of nodes are described in the comment of code.

#### Example

Following is an example of a network with two entries, one endpoint and two chains.

* Chain 1: entry1 -> channel1(with `30%` packet loss rate) -> restrict (1 pps, 1024 bps, buffer limited in 4096 bytes and 5 packets) -> endpoint
* Chain 2: entry2 -> channel2(with `10%` packet loss rate) -> endpoint

```go
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
	endpoint.Receive(func(packet base.Packet, now time.Time) []base.Event {
		if packet != nil {
			count.Inc()
			println("receive packet at", now.String())
			println("total", count.Load(), "packets received")
		}
		return nil
	})
	total := 20
	events := make([]base.Event, 0, total*2)
	for i := 0; i < 20; i++ {
		events = append(events, entry1.Send(base.RawPacket([]byte{0x01, 0x02})))
	}
	for i := 0; i < 20; i++ {
		events = append(events, entry2.Send(base.RawPacket([]byte{0x01, 0x02})))
	}
	event, cancel := base.NewPeriodicEvent(func(t time.Time) []base.Event {
		println("current time", t.String())
		return nil
	}, time.Second, time.Now())
	events = append(events, event)
	network.Start(events...)
	defer network.Stop()
	time.Sleep(time.Second * 30)
	cancel()
}
```

## Design

#### Architecture

The simulator is event driven, each event will be handled at the given time point, and generate subsequent events. Behaviors of nodes will be wrapped as events.

#### Event Loop

The event loop maintains a thread-local heap, in order to sort the events.

The loop is separated into two parts: fetch and drain.

* fetch: The main loop clear the event buffer, and put these events into the heap.
* drain: The main loop drain the heap until the heap only contains events at time point after the current time.

~~By now, the main loop lock a single os thread, but in the future, the main loop may run on a fork join pool.~~

~~Parallelized main loop is already implemented, main loop will split once the packet heap reach a given threshold, and exit after spinning a fixed rounds without any task. The active main loop will always exist but no more than a given limit.~~ (deprecated since it's hard to guarantee order of events, should be solved by separate multiple events heap)

## Contribution

#### Future work

* ~~parallelize main loop~~ (done)
* ~~implement commonly used protocol stack as a new node type~~ (will be implemented as different packet type)
* ~~separate send and pass to avoid cumulative error~~ (done)
* ~~Buffer overflow determination of *restrict node* should have a more accurate way~~ (done)
* split event heap when size of heap is large enough
* implement packets of commonly used protocol

#### Contributors

<a href="https://github.com/bytedance/ns-x/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=bytedance/ns-x"  alt="contributors"/>
</a>

Made with [contributors-img](https://contrib.rocks).