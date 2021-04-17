// Author: Enes Uysal
package main

type tree struct {
	root *node
}

type node struct {
	value    interface{}
	ty       int
	children []*node
}

func (T *tree) insert(parent *node, child *node) {

	if parent == nil {
		T.root = child
		return
	}

	child.children = make([]*node, 0)
	if parent.children != nil {
		parent.children = append(parent.children, child)
	}
}

func createNode(value interface{}, ty int, children []*node) *node {
	n := &node{
		value:    value,
		ty:       ty,
		children: children,
	}
	return n
}

func (T *tree) getNodeValue(path ...interface{}) *node {
	root := T.root

	for _, x := range path {
		if res, ok := x.(string); ok {
			if root.ty == 2 {
				for _, y := range root.value.([]*node) {
					if y.value == res {
						root = y
					}
				}
			} else {
				for _, y := range root.children {
					if y.value == res {
						root = y
					}
				}
			}
		} else if res, ok := x.(int); ok {
			if root.ty == 0 {
				root = root.children[0]
			}

			for i, x := range root.value.([]interface{}) {
				if i == res {
					root = x.(*node)
					if len(root.value.([]*node)) == 1 {
						root = root.value.([]*node)[0]
					}
				}
			}
		}
	}

	if root.ty == 0 && root.children != nil && len(root.children) == 1 {
		return root.children[0]
	}

	return root
}

func (T *tree) setNodeValue(value interface{}, path ...interface{}) {
	node := T.getNodeValue(path...)
	node.value = value
}
