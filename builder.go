package byte_ns

import (
	"byte-ns/base"
	"reflect"
	"strconv"
	"strings"
)

// Builder is a convenience tool to describe the whole network and build
type Builder interface {
	// Chain save the current chain and begin to describe a new chain
	Chain() Builder
	// Node connect the given node to the end of current chain, and name it if given node's name is not empty
	Node(node base.Node) Builder
	// Group connect the given nodes to the end of current chain one by one in order
	Group(nodes ...base.Node) Builder
	// NodeWithName same to Node, but use the given name instead of node's name
	NodeWithName(name string, node base.Node) Builder
	// GroupWithName same to Group, but name the whole group with the given name
	GroupWithName(name string, nodes ...base.Node) Builder
	// NodeByName find the node with the given name, and then connect it to the end of the chain
	NodeByName(name string) Builder
	// GroupByName find the group with the given name, and then perform the NodeGroup operation on it
	GroupByName(name string) Builder
	// Build actually connect the nodes with relation described before, any connection outside the builder will be overwritten
	// parameters are used to configure the network, return the built network, and a map from name to named nodes
	Build(loopLimit, emptySpinLimit, splitThreshold int) (*Network, map[string]base.Node)
}

type builder struct {
	nodes       map[base.Node]int
	names       map[string]base.Node
	groups      map[string][]base.Node
	current     base.Node
	connections map[base.Node]map[base.Node]interface{}
}

func NewBuilder() Builder {
	return &builder{
		nodes:       map[base.Node]int{},
		names:       map[string]base.Node{},
		groups:      map[string][]base.Node{},
		connections: map[base.Node]map[base.Node]interface{}{},
	}
}

func (b *builder) Chain() Builder {
	b.current = nil
	return b
}

func (b *builder) Node(node base.Node) Builder {
	return b.NodeWithName(node.Name(), node)
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
	if _, ok := b.nodes[node]; !ok {
		b.nodes[node] = len(b.nodes)
	}
	b.current = node
	if name != "" {
		b.names[name] = node
	}
	return b
}

func (b *builder) Group(nodes ...base.Node) Builder {
	return b.GroupWithName("", nodes...)
}

func (b *builder) GroupWithName(name string, nodes ...base.Node) Builder {
	if name != "" {
		b.groups[name] = nodes
	}
	builder := Builder(b)
	for _, node := range nodes {
		builder = builder.Node(node)
	}
	return builder
}

func (b *builder) NodeByName(name string) Builder {
	node, ok := b.names[name]
	if !ok {
		panic("no node with name: " + name)
	}
	return b.Node(node)
}

func (b *builder) GroupByName(name string) Builder {
	group, ok := b.groups[name]
	if !ok {
		panic("no group with name: " + name)
	}
	return b.Group(group...)
}

func (b *builder) Build(loopLimit, emptySpinLimit, splitThreshold int) (*Network, map[string]base.Node) {
	nodes := make([]base.Node, len(b.nodes))
	println("network summary: ")
	for node, index := range b.nodes {
		nodes[index] = node
	}
	for node, connection := range b.connections {
		node.SetNext(normalize(connection)...)
	}
	for index, node := range nodes {
		println(b.toString(node, index))
	}
	println()
	return NewNetwork(nodes, loopLimit, emptySpinLimit, splitThreshold), b.names
}

func (b *builder) toString(node base.Node, index int) string {
	sb := strings.Builder{}
	sb.WriteString("node ")
	sb.WriteString(strconv.Itoa(index))
	sb.WriteString(": {name: \"")
	sb.WriteString(node.Name())
	sb.WriteString("\", type: ")
	t := reflect.TypeOf(node)
	if t.Kind() == reflect.Ptr {
		sb.WriteString(t.Elem().Name())
	} else {
		sb.WriteString(t.Name())
	}
	sb.WriteString(", next: [")
	for _, n := range node.GetNext() {
		sb.WriteString(strconv.Itoa(b.nodes[n]))
		sb.WriteString(", ")
	}
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
