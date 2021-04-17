// Author: Enes Uysal
package yamlParser

type Tree struct {
	root *Node
}

type Node struct {
	Value    interface{}
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
		Value:    value,
		ty:       ty,
		children: children,
	}
	return n
}

func (T *Tree) GetNodeValue(path ...interface{}) *Node {
	root := T.root

	for _, x := range path {
		if res, ok := x.(string); ok {
			if root.ty == 2 {
				for _, y := range root.Value.([]*Node) {
					if y.Value == res {
						root = y
					}
				}
			} else {
				for _, y := range root.children {
					if y.Value == res {
						root = y
					}
				}
			}
		} else if res, ok := x.(int); ok {
			if root.ty == 0 {
				root = root.children[0]
			}

			for i, x := range root.Value.([]interface{}) {
				if i == res {
					root = x.(*Node)
					if len(root.Value.([]*Node)) == 1 {
						root = root.Value.([]*Node)[0]
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

func (T *Tree) SetNodeValue(value interface{}, path ...interface{}) {
	node := T.GetNodeValue(path...)
	node.Value = value
}
