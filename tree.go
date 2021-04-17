// Author: Enes Uysal
package yamlParser

type Tree struct {
	root *Node
}

type Node struct {
	value    interface{}
	ty       int
	children []*Node
}

func (T *Tree) insert(parent *Node, child *Node) {

	if parent == nil {
		T.root = child
		return
	}

	child.children = make([]*Node, 0)
	if parent.children != nil {
		parent.children = append(parent.children, child)
	}
}

func CreateNode(value interface{}, ty int, children []*Node) *Node {
	n := &Node{
		value:    value,
		ty:       ty,
		children: children,
	}
	return n
}

func (T *Tree) getNodeValue(path ...interface{}) *Node {
	root := T.root

	for _, x := range path {
		if res, ok := x.(string); ok {
			if root.ty == 2 {
				for _, y := range root.value.([]*Node) {
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
					root = x.(*Node)
					if len(root.value.([]*Node)) == 1 {
						root = root.value.([]*Node)[0]
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

func (T *Tree) setNodeValue(value interface{}, path ...interface{}) {
	node := T.getNodeValue(path...)
	node.value = value
}
