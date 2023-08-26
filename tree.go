package goad

import (
	"strings"
)

type Node struct {
	children []*Node
	handler  func(ctx *Context)
	param    string
	parent   *Node
}

func NewNode() *Node {
	return &Node{
		children: []*Node{},
		param:    "",
	}
}

func hasColonPrefix(param string) bool {
	return strings.HasPrefix(param, ":")
}

func (n *Node) Insert(path string, handler func(ctx *Context)) {
	node := n
	params := strings.Split(path, "/")
	for _, param := range params {
		child := node.findChild(param)

		if child == nil {
			child = &Node{
				param:    param,
				children: []*Node{},
				parent:   node,
			}
			node.children = append(node.children, child)
		}
		node = child
	}

	node.handler = handler
}

func (n *Node) findChild(param string) *Node {
	for _, child := range n.children {
		if child.param == param {
			return child
		}
	}

	return nil
}

func (n *Node) Search(path string) *Node {
	params := strings.Split(path, "/")
	return dfs(n, params)
}

func dfs(node *Node, params []string) *Node {
	currentParam := params[0]       //現在のparam
	isLastParam := len(params) == 1 //最後のparamかどうか

	for _, child := range node.children {
		if isLastParam {
			if hasColonPrefix(child.param) {
				return child
			}

			if child.param == currentParam {
				return child
			}

			continue
		}
		if !hasColonPrefix(child.param) && child.param != currentParam {
			continue
		}

		result := dfs(child, params[1:])

		if result != nil {
			return result
		}
	}
	return nil
}

func (n *Node) ParseParams(path string) map[string]string {
	node := n
	path = strings.TrimSuffix(path, "/")
	params := strings.Split(path, "/")
	paramMap := make(map[string]string)
	for i := len(params) - 1; i >= 0; i-- {
		if hasColonPrefix(node.param) {
			paramMap[node.param] = params[i]
		}
		node = node.parent
	}

	return paramMap
}
