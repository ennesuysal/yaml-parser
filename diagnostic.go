package main

import (
	"fmt"
	"regexp"
	"strings"
)

type lineType interface{}
type continuingLine struct{}
type arrayContinuing struct{}
type singleLine struct{}
type arraySingle struct{}
type arrayElement struct{}
type continuingString struct{}
type arrContStr struct{}
type continuingArr struct{}
type array struct{}

type diagnostic struct {
	root                 [][]interface{}
	tree                 *tree
	continuingStr        lineType
	continuingStrRoot    *node
	continuingStrIndent  float32
	continuingArr        lineType
	continuingArrIndent  float32
	continuingArrDim     int
	continuingArrFlag    int
	lastContString       lineType
	lastContStringIndent float32
	buffer               []string
}

func newYamlParser() *diagnostic {
	d := new(diagnostic)
	d.continuingStr = nil
	d.continuingArr = nil
	d.lastContString = nil
	d.continuingArrDim = 0

	rn := createNode("Root", 0, make([]*node, 0))
	d.tree = new(tree)
	d.tree.root = rn

	tmp := []interface{}{rn, float32(-1)}
	root := [][]interface{}{tmp}
	d.root = root

	return d
}

func analyze(line string) lineType {
	single, _ := regexp.MatchString(singleLineRgx, line)
	continuing, _ := regexp.MatchString(continuingLineRgx, line)
	arrayEl, _ := regexp.MatchString(arrayElementRgx, line)
	contStr, _ := regexp.MatchString(continuingStringRgx, line)
	contArr, _ := regexp.MatchString(continuingArrRgx, line)

	if arrayEl && contStr {
		return arrContStr{}
	} else if contStr {
		return continuingString{}
	} else if arrayEl && single {
		return arraySingle{}
	} else if arrayEl && continuing {
		return arrayContinuing{}
	} else if arrayEl {
		return arrayElement{}
	} else if single {
		return singleLine{}
	} else if continuing {
		return continuingLine{}
	} else if contArr {
		return continuingArr{}
	}

	return nil
}

func checkArray(p *node) *node {
	if p.ty == 1 {
		return p
	}

	for _, x := range p.children {
		if x.ty == 1 {
			return x
		}
	}
	return nil
}

func (d *diagnostic) writeBuffer() {
	if d.continuingStr == (continuingString{}) {
		d.continuingStr = nil
		line := strings.Join(d.buffer, "\n")
		n := createNode(line, 0, nil)
		d.tree.insert(d.continuingStrRoot, n)
		d.buffer = d.buffer[len(d.buffer):]
	}
}

func (d *diagnostic) scan(line string, indent float32) {
	ty := analyze(line)
	if ty == (arrayElement{}) {
		indent++
	}
	for lastRoot := d.root[len(d.root)-1][1].(float32); indent <= lastRoot; lastRoot = d.root[len(d.root)-1][1].(float32) {
		d.root = d.root[:len(d.root)-1]
	}

	switch ty.(type) {
	case arrayContinuing:
		d.writeBuffer()
		d.parseArrayCont(line, indent)
		if d.continuingArr == nil {
			return
		}

	case arraySingle:
		d.writeBuffer()
		if d.continuingArr == nil {
			d.parseArraySingle(line, indent)
			return
		}

	case arrayElement:
		d.writeBuffer()
		if d.continuingArr == nil {
			d.parseArrayElement(line, indent)
			return
		}

	case singleLine:
		d.writeBuffer()
		d.parseSingleLine(line, indent)
		if d.continuingArr == nil {
			return
		}

	case continuingLine:
		d.writeBuffer()
		d.parseContinuingLine(line, indent)
		if d.continuingArr == nil {
			return
		}

	case continuingString:
		d.continuingStr = continuingString{}
		d.continuingStrIndent = indent
		d.parseContinuingLine(parseContStr(line).(string)+":", indent)
		d.continuingStrRoot = d.root[len(d.root)-1][0].(*node)
		d.lastContString = nil
		if d.continuingArr == nil {
			return
		}

	case arrContStr:
		d.continuingStr = continuingString{}
		d.continuingStrIndent = indent
		d.parseArrayCont(parseArrContStr(line)+":", indent)
		d.continuingStrRoot = d.root[len(d.root)-1][0].(*node)
		d.lastContString = nil
		if d.continuingArr == nil {
			return
		}

	case continuingArr:
		d.lastContString = nil
		if d.continuingArr == nil {
			d.continuingArr = continuingArr{}
			d.continuingArrIndent = indent
		}
		d.continuingArrDim += 1
		return

	}

	if d.continuingStr != nil {
		d.buffer = append(d.buffer, line)
	}
	if d.continuingArr != nil {
		head := ""
		for i := 0; i < d.continuingArrDim; i++ {
			head += "- "
		}

		d.parseArrayElement(head+line, d.continuingArrIndent)
		fmt.Printf("%f\n", d.continuingArrIndent)
		d.continuingArr = nil
		d.continuingArrDim = 0
		d.continuingArrFlag = 0
	}
}
