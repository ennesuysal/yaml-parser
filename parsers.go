package main

import (
	"fmt"
	"regexp"
	"strings"
)

func (d *diagnostic) parseContinuingLine(line string, indent int) {
	out, _ := rgxShortcut(continuingLineRgx, line)
	if out != nil {
		fmt.Printf("ContLine: %v\n", out)
		key := createNode(out[0][1], 0, nil)
		d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		d.root = append(d.root, []interface{}{key, indent})
	}
}

func (d *diagnostic) parseSingleLine(line string) {
	out, _ := rgxShortcut(singleLineRgx, line)
	if out != nil {
		//fmt.Printf("Single Line: %v\n", out)
		key := createNode(out[0][1], 0, nil)
		value := createNode(out[0][2], 0, nil)

		d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		d.tree.insert(key, value)
	}
}

func (d *diagnostic) parseArrayElement(line string, indent int) {
	out, _ := rgxShortcut(arrayElementRgx, line)
	if out != nil {
		arrCount := strings.Count(out[0][1], "-")
		//spaceCount := strings.Count(out[0][1], " ")
		// var SpaceForOne int = spaceCount /arrCount

		//fmt.Printf("ArrayEl: %v\n", out)
		key := createNode(out[0][2], 0, make([]*node, 0))

		pa := checkArray(d.root[len(d.root)-1][0].(*node))

		if pa == nil {
			pa = createNode(make([]interface{}, 0), 1, nil)
			d.tree.insert(d.root[len(d.root)-1][0].(*node), pa)
			d.root = append(d.root, []interface{}{pa, indent})
			fmt.Printf("New Array Indent: %d\n", indent)
		}
		var condIndent int
		i := 0
		for ; i < arrCount-1; i++ {
			tmp := createNode(make([]interface{}, 0), 1, nil)
			pa.value = append(pa.value.([]interface{}), tmp)
			pa = tmp
			condIndent = indent + (i+1)*2
			d.root = append(d.root, []interface{}{pa, condIndent})
			fmt.Printf("Array Indent: %d\n", condIndent)
		}

		//fmt.Printf("%v", pa.value)
		pa.value = append(pa.value.([]interface{}), key)

		if out[0][3] != "" {
			child := createNode(out[0][3], 0, nil)
			d.tree.insert(key, child)
		} else {
			d.root = append(d.root, []interface{}{key, indent + i*2 + 1})
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
