package main

import (
	"fmt"
	"regexp"
	"strings"
)

type lineType interface{}
type continuingLine struct{}
type arrayAndContinuing struct{}
type singleLine struct{}
type arrayAndSingle struct{}
type arrayElement struct{}
type continuingString struct{}
type arrAndContStr struct{}
type continuingArr struct{}
type singleLineArr struct{}
type arrArr struct{}
type arrSLAE struct{}

type diagnostic struct {
	root                    [][]interface{}
	tree                    *tree
	continuingStr           lineType
	continuingStrRoot       *node
	continuingStrIndent     float32
	continuingArr           lineType
	continuingArrIndent     float32
	continuingArrLastIndent float32
	continuingArrSpaceCount float32
	continuingArrDim        int
	continuingArrFlag       int
	lastContString          lineType
	lastContStringIndent    float32
	buffer                  []string
}

func newYamlParser() *diagnostic {
	d := new(diagnostic)
	d.continuingStr = nil
	d.continuingArr = nil
	d.lastContString = nil
	d.continuingArrDim = 0
	d.continuingArrSpaceCount = 0
	d.continuingArrLastIndent = 0

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
	arrArr_, _ := regexp.MatchString(arrRgx, line)
	SLAE_, _ := regexp.MatchString(singleLineArrayRgx, line)

	if arrayEl && contStr {
		return arrAndContStr{}
	} else if contStr {
		return continuingString{}
	} else if SLAE_ && arrayEl {
		return arrSLAE{}
	} else if arrayEl && single {
		return arrayAndSingle{}
	} else if arrayEl && continuing {
		return arrayAndContinuing{}
	} else if arrayEl && arrArr_ {
		return arrArr{}
	} else if arrayEl {
		return arrayElement{}
	} else if SLAE_ {
		return singleLineArr{}
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

	if d.continuingArr != nil && ty != (continuingArr{}) {
		head := ""
		for i := 0; i < d.continuingArrDim; i++ {
			head += "-"
			for j := float32(0); j < d.continuingArrSpaceCount; j++ {
				head += " "
			}
		}
		line = head + line
		indent = d.continuingArrIndent
		fmt.Printf("%s\n", line)
		d.continuingArr = nil
		d.continuingArrDim = 0
		d.continuingArrLastIndent = 0
		d.continuingArrSpaceCount = 0
		d.continuingArrFlag = 0
	}

	switch ty.(type) {
	case arrayAndContinuing:
		d.writeBuffer()
		if d.continuingArr == nil {
			d.parseArrayCont(line, indent)
			return
		}

	case arrayAndSingle:
		d.writeBuffer()
		if d.continuingArr == nil {
			d.parseArraySingle(line, indent)
			return
		}

	case arrayElement:
		d.writeBuffer()
		if d.continuingArr == nil {
			d.parseArrayElement(line, indent, false)
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

	case arrAndContStr:
		if d.continuingArr == nil {
			d.parseArrayCont(line, indent)
		}

	case continuingArr:
		d.lastContString = nil
		d.continuingArrSpaceCount = (indent-d.continuingArrLastIndent)/2 - 1
		d.continuingArrLastIndent = indent
		if d.continuingArr == nil {
			d.continuingArr = continuingArr{}
			d.continuingArrIndent = indent
		}
		d.continuingArrDim += 1
		return

	case singleLineArr:
		d.parseSingleLineArray(line, indent)

	case arrArr:
		d.parseSLAE(line, indent)

	case arrSLAE:
		d.parseArrSLAE(line, indent)
	}

	if d.continuingStr != nil {
		d.buffer = append(d.buffer, line)
	}
}
