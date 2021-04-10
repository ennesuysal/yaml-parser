package main

import (
	"math"
	"regexp"
	"strings"
)

func (d *diagnostic) parseContinuingLine(line string, indent float32) {
	out, _ := rgxShortcut(continuingLineRgx, line)
	if out != nil {
		key := createNode(out[0][1], 0, make([]*node, 0))
		if d.lastContString == (arrayElement{}) && d.lastContStringIndent < indent && d.root[len(d.root)-1][0].(*node).ty == 1 {
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
		if d.lastContString == (arrayElement{}) && d.lastContStringIndent < indent && d.root[len(d.root)-1][0].(*node).ty == 1 {
			r := d.root[len(d.root)-1][0].(*node).value.([]interface{})
			r[len(r)-1] = append(r[len(r)-1].([]*node), key)
		} else {
			d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		}

		d.tree.insert(key, value)
	}
}

func (d *diagnostic) parseArraySingle(line string, indent float32) {
	out, _ := rgxShortcut(singleLineRgx, line)
	child := createNode(out[0][2], 0, nil)
	key, _ := d.parseArrayElement(line, indent)
	d.tree.insert(key, child)
}

func (d *diagnostic) parseArrayCont(line string, indent float32) {
	key, cIndent := d.parseArrayElement(line, indent)
	d.root = append(d.root, []interface{}{key, cIndent})
}

func (d *diagnostic) parseArrayElement(line string, indent float32) (*node, float32) {
	out, _ := rgxShortcut(arrayElementRgx, line)
	if out != nil {
		arrCount := float32(strings.Count(out[0][1], "-"))
		spcCount := float32(strings.Count(out[0][1], " "))

		spcPerArr := int(math.Ceil(float64(spcCount / arrCount)))

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
			condIndent += 2 + float32(spcPerArr)*2
			d.root = append(d.root, []interface{}{pa, condIndent - 0.5})
		}

		nodeArray := make([]*node, 0)
		nodeArray = append(nodeArray, key)
		pa.value = append(pa.value.([]interface{}), nodeArray)

		d.lastContString = arrayElement{}
		d.lastContStringIndent = condIndent
		return key, condIndent + 2 + float32(spcPerArr)*2
	}
	return nil, -2
}

func parseContStr(line string) interface{} {
	re, _ := regexp.Compile(continuingStringRgx)
	out := re.FindStringSubmatch(line)
	if out != nil {
		return out[1]
	}
	return nil
}

func parseArrContStr(line string) string {
	re, _ := regexp.Compile(arrayElementRgx)
	out := re.FindStringSubmatch(line)
	if out != nil {
		return out[1] + out[2]
	}
	return ""
}
