package byte_ns

import (
	"reflect"
	"strconv"
	"strings"
)

type Builder interface {
	Node(node Node) Builder
	Group(nodes ...Node) Builder
	NodeWithName(name string, node Node) Builder
	GroupWithName(name string, nodes ...Node) Builder
	NodeByName(name string) Builder
	GroupByName(name string) Builder
	Chain() Builder
	Build(loopLimit, emptySpinLimit, splitThreshold int) (*Network, map[string]Node)
}

type builder struct {
	nodes       map[Node]int
	names       map[string]Node
	groups      map[string][]Node
	current     Node
	connections map[Node]map[Node]interface{}
}

func NewBuilder() Builder {
	return &builder{
		nodes:       map[Node]int{},
		names:       map[string]Node{},
		groups:      map[string][]Node{},
		connections: map[Node]map[Node]interface{}{},
	}
}

func (b *builder) Chain() Builder {
	b.current = nil
	return b
}

func (b *builder) Node(node Node) Builder {
	return b.NodeWithName(node.Name(), node)
}

func (b *builder) NodeWithName(name string, node Node) Builder {
	if b.current != nil {
		connection, ok := b.connections[b.current]
		if !ok {
			connection = map[Node]interface{}{}
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

func (b *builder) Group(nodes ...Node) Builder {
	return b.GroupWithName("", nodes...)
}

func (b *builder) GroupWithName(name string, nodes ...Node) Builder {
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

func (b *builder) Build(loopLimit, emptySpinLimit, splitThreshold int) (*Network, map[string]Node) {
	nodes := make([]Node, len(b.nodes))
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

func (b *builder) toString(node Node, index int) string {
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

func normalize(nodes map[Node]interface{}) []Node {
	result := make([]Node, 0, len(nodes))
	for node := range nodes {
		result = append(result, node)
	}
	return result
}
