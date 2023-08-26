package goad

import (
	"net/http"
	"strings"
)

type Node struct {
	children []*Node
	handler  func(w http.ResponseWriter, r *http.Request)
	param    string
}

func NewNode() Node {
	return Node{
		children: []*Node{},
		param:    "",
	}
}

func hasColonPrefix(param string) bool {
	return strings.HasPrefix(param, ":")
}

func (n *Node) Insert(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	node := n
	params := strings.Split(path, "/")
	for _, param := range params {
		child := node.findChild(param)

		if child == nil {
			child = &Node{
				param:    param,
				children: []*Node{},
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

func (n *Node) Search(path string) func(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(path, "/")
	result := dfs(n, params)

	if result == nil {
		return nil
	}

	return result.handler
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
