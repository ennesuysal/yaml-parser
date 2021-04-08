package main

import (
	"regexp"
	"strings"
)

func (d *diagnostic) parseContinuingLine(line string, indent float32) {
	out, _ := rgxShortcut(continuingLineRgx, line)
	if out != nil {
		key := createNode(out[0][1], 0, make([]*node, 0))
		if d.last == (arrayElement{}) && d.lastIndent < indent && d.root[len(d.root)-1][0].(*node).ty == 1 {
			r := d.root[len(d.root)-1][0].(*node).value.([]interface{})
			r[len(r)-1] = append(r[len(r)-1].([]*node), key)
		} else {
			d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		}
		d.root = append(d.root, []interface{}{key, indent})
	}
}

func (d *diagnostic) parseSingleLine(line string, indent float32) {
	out, _ := rgxShortcut(singleLineRgx, line)
	if out != nil {
		key := createNode(out[0][1], 0, make([]*node, 0))
		value := createNode(out[0][2], 0, nil)
		if d.last == (arrayElement{}) && d.lastIndent < indent && d.root[len(d.root)-1][0].(*node).ty == 1 {
			r := d.root[len(d.root)-1][0].(*node).value.([]interface{})
			r[len(r)-1] = append(r[len(r)-1].([]*node), key)
		} else {
			d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		}

		d.tree.insert(key, value)
	}
}

func (d *diagnostic) parseArrayElement(line string, indent float32) {
	out, _ := rgxShortcut(arrayElementRgx, line)
	if out != nil {
		arrCount := float32(strings.Count(out[0][1], "-"))

		key := createNode(out[0][2], 0, make([]*node, 0))

		pa := checkArray(d.root[len(d.root)-1][0].(*node))

		if pa == nil {
			pa = createNode(make([]interface{}, 0), 1, nil)
			d.tree.insert(d.root[len(d.root)-1][0].(*node), pa)
			d.root = append(d.root, []interface{}{pa, indent - 0.5})
		}

		condIndent := indent
		i := float32(0)
		for ; i < arrCount-1; i++ {
			tmp := createNode(make([]interface{}, 0), 1, nil)
			pa.value = append(pa.value.([]interface{}), tmp)
			pa = tmp
			condIndent += (i + 1) * 2
			d.root = append(d.root, []interface{}{pa, condIndent})
		}
		nodeArray := make([]*node, 0)
		nodeArray = append(nodeArray, key)
		pa.value = append(pa.value.([]interface{}), nodeArray)

		d.last = arrayElement{}
		d.lastIndent = condIndent

		if out[0][3] != "" {
			child := createNode(out[0][3], 0, nil)
			d.tree.insert(key, child)
		} else {
			d.root = append(d.root, []interface{}{key, indent + (i+1)*2})
		}
	}
}

func parseContStr(line string) interface{} {
	re, _ := regexp.Compile(continuingStringRgx)
	out := re.FindStringSubmatch(line)
	if out != nil {
		return out[1]
	}
	return nil
}
