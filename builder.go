package ns_x

import (
	"github.com/bytedance/ns-x/v2/base"
	"reflect"
	"strconv"
	"strings"
)

// Builder is a convenience tool to describe the whole network and build
// when a node is described for the first time, a unique id is assigned to it, which will be the index of node in the built network
type Builder interface {
	// Chain save the current chain and begin to describe a new chain
	Chain() Builder
	// Node connect the given node to the end of current chain
	Node(node base.Node) Builder
	// Group insert a group to current chain, which means end of current chain will be connected to in node, and the end of current chain will be set to out node
	Group(inName, outName string) Builder
	// NodeWithName same to Node, but name it with the given name
	NodeWithName(name string, node base.Node) Builder
	// GroupWithName same to Group, but name the whole group with the given name
	GroupWithName(name string, inName, outName string) Builder
	// NodeOfName find the node with the given name, and then connect it to the end of the chain
	NodeOfName(name string) Builder
	// GroupOfName find the group with the given name, and then perform the Group operation on it
	GroupOfName(name string) Builder
	// Summary print the structure of the network to standard output
	Summary() Builder
	// Build actually connect the nodes with relation described before, any connection outside the builder will be overwritten
	// parameters are used to configure the network, return the built network, and a map from name to named nodes
	Build() (*Network, map[string]base.Node)
}

type group struct {
	inName, outName string
}

type builder struct {
	nodeToID    map[base.Node]int
	nameToNode  map[string]base.Node
	nodeToName  map[base.Node]string
	nameToGroup map[string]*group
	current     base.Node
	connections map[base.Node]map[base.Node]interface{}
}

func NewBuilder() Builder {
	return &builder{
		nodeToID:    map[base.Node]int{},
		nameToNode:  map[string]base.Node{},
		nodeToName:  map[base.Node]string{},
		nameToGroup: map[string]*group{},
		connections: map[base.Node]map[base.Node]interface{}{},
	}
}

func (b *builder) Chain() Builder {
	b.current = nil
	return b
}

func (b *builder) Node(node base.Node) Builder {
	return b.NodeWithName("", node)
}

func (b *builder) NodeWithName(name string, node base.Node) Builder {
	if b.current != nil {
		connection, ok := b.connections[b.current]
		if !ok {
			connection = map[base.Node]interface{}{}
			b.connections[b.current] = connection
		}
		connection[node] = nil
	}
	if _, ok := b.nodeToID[node]; !ok {
		b.nodeToID[node] = len(b.nodeToID)
	}
	if name != "" {
		b.nameToNode[name] = node
		b.nodeToName[node] = name
	}
	b.current = node
	return b
}

func (b *builder) Group(inName, outName string) Builder {
	return b.GroupWithName("", inName, outName)
}

func (b *builder) GroupWithName(name string, inName, outName string) Builder {
	if name != "" {
		b.nameToGroup[name] = &group{inName: inName, outName: outName}
	}
	in := b.requireNodeByName(inName)
	out := b.requireNodeByName(outName)
	b.Node(in)
	b.current = out
	return b
}

func (b *builder) NodeOfName(name string) Builder {
	return b.Node(b.requireNodeByName(name))
}

func (b *builder) GroupOfName(name string) Builder {
	group, ok := b.nameToGroup[name]
	if !ok {
		panic("no group with name: " + name)
	}
	return b.Group(group.inName, group.outName)
}

func (b *builder) Summary() Builder {
	nodes := make([]base.Node, len(b.nodeToID))
	println("network summary: ")
	for node, index := range b.nodeToID {
		nodes[index] = node
	}
	for index, node := range nodes {
		println(b.toString(node, index))
	}
	println()
	return b
}

func (b *builder) Build() (*Network, map[string]base.Node) {
	nodes := make([]base.Node, len(b.nodeToID))
	for node, index := range b.nodeToID {
		nodes[index] = node
	}
	for node, connection := range b.connections {
		node.SetNext(normalize(connection)...)
	}
	return NewNetwork(nodes), b.nameToNode
}

func (b *builder) toString(node base.Node, index int) string {
	sb := strings.Builder{}
	sb.WriteString("node ")
	sb.WriteString(strconv.Itoa(index))
	sb.WriteString(": {name: \"")
	sb.WriteString(b.nodeToName[node])
	sb.WriteString("\", type: ")
	t := reflect.TypeOf(node)
	if t.Kind() == reflect.Ptr {
		sb.WriteString(t.Elem().Name())
	} else {
		sb.WriteString(t.Name())
	}
	sb.WriteString(", next: [")
	connection := b.connections[node]
	next := make([]string, 0, len(connection))
	for n := range connection {
		next = append(next, strconv.Itoa(b.nodeToID[n]))
	}
	sb.WriteString(strings.Join(next, ","))
	sb.WriteString("]}")
	return sb.String()
}

func normalize(nodes map[base.Node]interface{}) []base.Node {
	result := make([]base.Node, 0, len(nodes))
	for node := range nodes {
		result = append(result, node)
	}
	return result
}

func (b *builder) requireNodeByName(name string) base.Node {
	if name == "" {
		panic("name cannot be empty string")
	}
	node, ok := b.nameToNode[name]
	if !ok {
		panic("no node with name " + name)
	}
	return node
}
