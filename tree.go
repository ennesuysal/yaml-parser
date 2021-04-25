// Author: Enes Uysal
package yamlParser

import (
	"os"
)

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

			for i, x := range root.Value.([]*Node) {
				if i == res {
					root = x
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

func writeFileHelper(n *Node, txt string, indent int, parent int) string {
	if n.ty != 1 && n.ty != 2 && parent != 1 && len(n.children) == 0 {
		txt += " "
	} else if n.ty != 1 && n.ty != 2 {
		for i := 0; i < indent; i++ {
			txt += "  "
		}
	}

	if n.ty == 0 {
		if parent == 1 {
			txt += "- "
			indent++
		}
		txt += n.Value.(string)
		if len(n.children) > 0 {
			txt += ":"
			if !(len(n.children) == 1 && n.children[0].ty == 0 && len(n.children[0].children) == 0) {
				txt += "\n"
			}
			for i := 0; i < len(n.children); i++ {
				txt = writeFileHelper(n.children[i], txt, indent+1, 0)
			}
		} else {
			txt += "\n"
		}
	} else if n.ty == 1 {
		if parent == 1 {
			for i := 0; i < indent; i++ {
				txt += "  "
			}
			txt += "-\n"
			indent += 1
		}
		for _, x := range n.Value.([]*Node) {
			txt = writeFileHelper(x, txt, indent, 1)
		}
	} else if n.ty == 2 {
		for i, x := range n.Value.([]*Node) {
			if i == 0 {
				txt = writeFileHelper(x, txt, indent, 1)
			} else {
				txt = writeFileHelper(x, txt, indent+1, 0)
			}
		}
	}

	return txt
}

func (T *Tree) WriteFile(path string) {
	f, err := os.Create(path)
	if err == nil {
		var indent = 0
		root := T.root
		val := ""
		for _, x := range root.children {
			val += writeFileHelper(x, "", indent, 0)
		}
		f.Write([]byte(val))
		f.Close()
	}
}
