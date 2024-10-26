package gee

import "strings"

type trieNode struct {
	Pattern  string      // 待匹配路由，e.g. /p/:lang
	part     string      // 路由中间的一部分，比如:lang
	children []*trieNode // 子节点
	isWild   bool        // 是否精准匹配，当part中含有 ':' 或 '*' 时为 true
}

func (n *trieNode) matchChild(part string) *trieNode {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *trieNode) matchChildren(part string) []*trieNode {
	nodes := make([]*trieNode, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *trieNode) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.Pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &trieNode{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *trieNode) search(parts []string, height int) *trieNode {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.Pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
