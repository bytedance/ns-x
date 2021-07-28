package byte_ns

import (
	"reflect"
	"strconv"
	"strings"
)

type Builder interface {
	Node(node Node) Builder
	Nodes(nodes ...Node) Builder
	NodeByName(name string) Builder
	Chain() Builder
	Build() (*Network, map[string]Node)
}

type builder struct {
	nodes   map[Node]int
	names   map[string]Node
	current Node
}

func NewBuilder() Builder {
	return &builder{
		nodes: map[Node]int{},
		names: map[string]Node{},
	}
}

func (b *builder) Chain() Builder {
	b.current = nil
	return b
}

func (b *builder) Node(node Node) Builder {
	if b.current != nil {
		b.current.SetNext(append(b.current.GetNext(), node)...)
	}
	if _, ok := b.nodes[node]; !ok {
		b.nodes[node] = len(b.nodes)
	}
	b.current = node
	if name := node.Name(); name != "" {
		b.names[name] = node
	}
	return b
}

func (b *builder) Nodes(nodes ...Node) Builder {
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

func (b *builder) Build() (*Network, map[string]Node) {
	nodes := make([]Node, len(b.nodes))
	println("network summary: ")
	for node, index := range b.nodes {
		nodes[index] = node
	}
	for index, node := range nodes {
		println(b.toString(node, index))
	}
	println()
	return NewNetwork(nodes), b.names
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
