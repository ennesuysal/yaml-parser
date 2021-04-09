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

func (T *tree) travel(p *node) {
	if p == nil {
		return
	}

	if len(p.children) == 0 {
		if p.ty == 0 {
			//fmt.Printf("%s\n", p.value)
		} else {
			for _, x := range p.children {
				if p.ty == 0 {
					//		fmt.Printf("%s\n", p.value)
				} else {
					T.travel(x)
				}
			}
		}
		return
	}

	for _, x := range p.children {
		T.travel(x)
	}
}
