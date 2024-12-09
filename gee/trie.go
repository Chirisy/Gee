package gee

import "strings"

// 前缀树,实现动态路由
type node struct {
	pattern  string  //待匹配的路由,ex: /p/:lang
	part     string  //路由中的一部分,ex: :lang
	children []*node //子节点,比如 [doc, tutorial, intro]
	isWild   bool    //是否精确匹配,part含有:或*时为true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child //找到第一个匹配成功的子节点,返回
		}
	}
	return nil
}

// 匹配子节点,找到所有匹配成功的节点,返回
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child) //找到所有匹配成功的子节点,返回
		}
	}
	return nodes
}

// 插入路由节点,无匹配的节点就新建
func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1) //继续在下一层匹配
}

// 查询路由
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part) //匹配本层的节点
	//对每一个节点进行搜索
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
