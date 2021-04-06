package main

import (
	"errors"
	"regexp"
)

const (
	singleLineRgx     = `(.+)\s*:\s*(.+)`
	continuingLineRgx = `([^\s]+)\s*:\s*$`
	arrayElementRgx   = `-\s*([^\s]*)\s*:\s*([^\s]*)\s*`
)

type lineType interface{}
type continuingLine struct{}
type singleLine struct{}
type arrayElement struct{}
type array struct{}

type diagnostic struct {
	root [][]interface{}
}

func rgx_shortcut(rgx string, txt string) ([][][]byte, error) {
	r, err := regexp.Compile(rgx)

	if err != nil {
		return nil, errors.New("Invalid regex!")
	}

	if !r.Match([]byte(txt)) {
		return nil, errors.New("string does not match!")
	}
	match := r.FindAllSubmatch([]byte(txt), -1)
	return match, nil
}

func trim(line string) (int, string) {
	i := 0
	for ; i < len(line); i += 2 {
		if line[i] != ' ' && line[i] != '\t' {
			break
		}
	}

	indent := i

	if analyze(line) == (arrayElement{}) {
		indent += 1
	}

	return indent, line[i:]
}

func createNode(value interface{}, ty int, children []*node) *node {
	n := &node{
		value:    value,
		ty:       ty,
		children: children,
	}
	return n
}

func analyze(line string) lineType {
	single, _ := regexp.MatchString(singleLineRgx, line)
	continuing, _ := regexp.MatchString(continuingLineRgx, line)
	arrayEl, _ := regexp.MatchString(arrayElementRgx, line)

	if arrayEl {
		return arrayElement{}
	} else if single {
		return singleLine{}
	} else if continuing {
		return continuingLine{}
	}

	return nil
}

func checkArray(p *node) *node {
	for _, x := range p.children {
		if x.ty == 1 {
			return x
		}
	}
	return nil
}

func (d *diagnostic) scan(line string, indent int, T *tree) {
	for lastRoot := d.root[len(d.root)-1][1].(int); indent <= lastRoot; lastRoot = d.root[len(d.root)-1][1].(int) {
		d.root = d.root[:len(d.root)-1]
	}

	ty := analyze(line)
	switch ty.(type) {
	case arrayElement:
		out, _ := rgx_shortcut(arrayElementRgx, line)
		if out != nil {
			key := createNode(out[0][1], 0, make([]*node, 0))

			pa := checkArray(d.root[len(d.root)-1][0].(*node))

			if pa == nil {
				pa = createNode(make([]interface{}, 0), 1, nil)

				T.insert(d.root[len(d.root)-1][0].(*node), pa)
			}

			pa.value = append(pa.value.([]interface{}), key)

			if string(out[0][2]) != "" {
				child := createNode(out[0][2], 0, nil)
				T.insert(key, child)
			}

			d.root = append(d.root, []interface{}{pa, indent})
			d.root = append(d.root, []interface{}{key, indent})
		}
		break

	case singleLine:
		out, _ := rgx_shortcut(singleLineRgx, line)
		if out != nil {

			key := createNode(out[0][1], 0, nil)
			value := createNode(out[0][2], 0, nil)

			T.insert(d.root[len(d.root)-1][0].(*node), key)
			T.insert(key, value)
		}

	case continuingLine:
		out, _ := rgx_shortcut(continuingLineRgx, line)
		if out != nil {

			key := createNode(out[0][1], 0, nil)

			T.insert(d.root[len(d.root)-1][0].(*node), key)

			d.root = append(d.root, []interface{}{key, indent})
		}
	}
}
