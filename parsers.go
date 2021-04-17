// Author: Enes Uysal
package yamlParser

import (
	"math"
	"regexp"
	"strings"
)

func (d *Diagnostic) parseContinuingLine(line string, indent float32) {
	out, _ := rgxShortcut(continuingLineRgx, line)
	if out != nil {
		key := CreateNode(out[0][1], 0, make([]*Node, 0))
		if d.lastContString == (arrayElement{}) && d.lastContStringIndent < indent && d.root[len(d.root)-1][0].(*Node).ty == 1 {
			r := d.root[len(d.root)-1][0].(*Node).value.([]interface{})
			r[len(r)-1].(*Node).value = append(r[len(r)-1].(*Node).value.([]*Node), key)
		} else {
			d.tree.insert(d.root[len(d.root)-1][0].(*Node), key)
		}
		d.root = append(d.root, []interface{}{key, indent})
	}
}

func (d *Diagnostic) parseSingleLine(line string, indent float32) {
	out, _ := rgxShortcut(singleLineRgx, line)
	if out != nil {
		key := CreateNode(out[0][1], 0, make([]*Node, 0))
		value := CreateNode(out[0][2], 0, nil)
		if d.lastContString == (arrayElement{}) && d.lastContStringIndent < indent && d.root[len(d.root)-1][0].(*Node).ty == 1 {
			r := d.root[len(d.root)-1][0].(*Node).value.([]interface{})
			r[len(r)-1].(*Node).value = append(r[len(r)-1].(*Node).value.([]*Node), key)
		} else {
			d.tree.insert(d.root[len(d.root)-1][0].(*Node), key)
		}

		d.tree.insert(key, value)
	}
}

func (d *Diagnostic) parseArraySingle(line string, indent float32) {
	out, _ := rgxShortcut(singleLineRgx, line)
	child := CreateNode(out[0][2], 0, nil)
	key, _ := d.parseArrayElement(line, indent, false)
	d.tree.insert(key, child)
}

func (d *Diagnostic) parseArrayCont(line string, indent float32) {
	key, cIndent := d.parseArrayElement(line, indent, false)
	d.root = append(d.root, []interface{}{key, cIndent})
}

func (d *Diagnostic) parseArrayElement(line string, indent float32, arr bool) (*Node, float32) {
	out, _ := rgxShortcut(arrayElementRgx, line)
	if out != nil {
		arrCount := float32(strings.Count(out[0][1], "-"))
		spcCount := float32(strings.Count(out[0][1], " "))

		spcPerArr := int(math.Ceil(float64(spcCount / arrCount)))

		key := CreateNode(out[0][2], 0, make([]*Node, 0))

		pa := checkArray(d.root[len(d.root)-1][0].(*Node))

		if pa == nil {
			pa = CreateNode(make([]interface{}, 0), 1, nil)
			d.tree.insert(d.root[len(d.root)-1][0].(*Node), pa)
			d.root = append(d.root, []interface{}{pa, indent - 0.5})
		}

		condIndent := indent
		i := float32(0)
		for ; i < arrCount-1; i++ {
			tmp := CreateNode(make([]interface{}, 0), 1, nil)
			pa.value = append(pa.value.([]interface{}), tmp)
			pa = tmp
			condIndent += 2 + float32(spcPerArr)*2
			d.root = append(d.root, []interface{}{pa, condIndent - 0.5})
		}

		d.lastContString = arrayElement{}
		d.lastContStringIndent = condIndent

		if arr {
			return pa, condIndent + 2 + float32(spcPerArr)*2
		}

		nodeArray := CreateNode(make([]*Node, 0), 2, nil)
		nodeArray.value = append(nodeArray.value.([]*Node), key)
		pa.value = append(pa.value.([]interface{}), nodeArray)

		return key, condIndent + 2 + float32(spcPerArr)*2
	}
	return nil, -2
}

func (d *Diagnostic) parseSingleLineArray(line string, indent float32) {
	out, _ := rgxShortcut(singleLineArrayRgx, line)
	if out != nil {
		key := CreateNode(out[0][1], 0, make([]*Node, 0))
		value := sLineArrayHelper(out[0][2])
		if d.lastContString == (arrayElement{}) && d.lastContStringIndent < indent && d.root[len(d.root)-1][0].(*Node).ty == 1 {
			r := d.root[len(d.root)-1][0].(*Node).value.([]interface{})
			r[len(r)-1].(*Node).value = append(r[len(r)-1].(*Node).value.([]*Node), key)
		} else {
			d.tree.insert(d.root[len(d.root)-1][0].(*Node), key)
		}
		d.tree.insert(key, value)
	}
}

func (d *Diagnostic) parseSLAE(line string, indent float32) {
	out, _ := rgxShortcut(arrRgx, line)
	child := sLineArrayHelper(out[0][1])
	key, _ := d.parseArrayElement(line, indent, true)
	key.value = append(key.value.([]interface{}), child)
}

func (d *Diagnostic) parseArrSLAE(line string, indent float32) {
	out, _ := rgxShortcut(singleLineArrayRgx, line)
	child := sLineArrayHelper(out[0][2])
	key, _ := d.parseArrayElement(line, indent, false)
	d.tree.insert(key, child)
}

func sLineArrayHelper(line string) *Node {
	queue := make([]interface{}, 0)
	buffer := ""
	for _, x := range line {
		if x == '[' {
			arr := CreateNode(make([]interface{}, 0), 1, nil)
			queue = append(queue, arr)
		} else if x == ']' {
			if buffer != "" {
				n := CreateNode(buffer, 0, nil)
				nodeArray := make([]*Node, 0)
				nodeArray = append(nodeArray, n)
				queue[len(queue)-1].(*Node).value = append(queue[len(queue)-1].(*Node).value.([]interface{}), nodeArray)
				buffer = ""
			}
			if len(queue) > 1 {
				queue[len(queue)-2].(*Node).value = append(queue[len(queue)-2].(*Node).value.([]interface{}), queue[len(queue)-1])
				queue = queue[:len(queue)-1]
			} else {
				break
			}
		} else if x == ',' {
			if buffer != "" {
				n := CreateNode(buffer, 0, nil)
				nodeArray := make([]*Node, 0)
				nodeArray = append(nodeArray, n)
				queue[len(queue)-1].(*Node).value = append(queue[len(queue)-1].(*Node).value.([]interface{}), nodeArray)
				buffer = ""
			}
		} else {
			if x != ' ' {
				buffer += string(x)
			}
		}
	}

	return queue[0].(*Node)
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
