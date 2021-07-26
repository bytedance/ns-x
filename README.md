## network simulator

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

### Design

#### Main Loop

The main loop maintains a thread-local heap, in order to sort the packets.

The loop is separated into two parts: fetch and drain.

* fetch: The main loop clear the packet buffer of all the nodes, and put these packets into the heap.
* drain: The main loop drain the heap until the heap only contains packets with emit time after the current time.

By now, the main loop lock a single os thread, but in the future, the main loop may run on a fork join pool.
