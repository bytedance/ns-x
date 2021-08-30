package node

import "github.com/bytedance/ns-x/v2/base"

// Option is applied on a node to make some changes
type Option func(node base.Node)

// apply given options one by one on the given node
func apply(node base.Node, options ...Option) {
	for _, option := range options {
		option(node)
	}
}
