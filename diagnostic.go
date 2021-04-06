package main

import (
	"regexp"
	"strings"
)

const (
	singleLineRgx       = `(.+)\s*:\s*([^\|>]+)$`
	continuingLineRgx   = `([^\s]+)\s*:\s*$`
	arrayElementRgx     = `-\s*([^\s]*)\s*:\s*([^\s]*)\s*`
	continuingStringRgx = `(.+)\s*:\s*[>\|]\s*$`
	continuingArrRgx    = `^\s*-\s*$`
)

type lineType interface{}
type continuingLine struct{}
type singleLine struct{}
type arrayElement struct{}
type continuingString struct{}
type continuingArr struct{}
type array struct{}

type diagnostic struct {
	root             [][]interface{}
	tree             *tree
	continuing       lineType
	continuingRoot   *node
	continuingIndent int
	buffer           []string
}

func newYamlParser() *diagnostic {
	d := new(diagnostic)
	d.continuing = nil

	rn := createNode("Root", 0, make([]*node, 0))
	d.tree = new(tree)
	d.tree.root = rn

	tmp := []interface{}{rn, -1}
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

	if arrayEl {
		return arrayElement{}
	} else if single {
		return singleLine{}
	} else if continuing {
		return continuingLine{}
	} else if contStr {
		return continuingString{}
	} else if contArr {
		return continuingArr{}
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

func (d *diagnostic) writeBuffer() {
	if d.continuing != nil {
		d.continuing = nil
		line := strings.Join(d.buffer, "\n")
		n := createNode(line, 0, nil)
		d.tree.insert(d.continuingRoot, n)
		d.buffer = d.buffer[:len(d.buffer)]
	}
}

func (d *diagnostic) scan(line string, indent int) {
	for lastRoot := d.root[len(d.root)-1][1].(int); indent <= lastRoot; lastRoot = d.root[len(d.root)-1][1].(int) {
		d.root = d.root[:len(d.root)-1]
	}

	if d.continuing != nil {
		switch d.continuing.(type) {
		case continuingString:
			d.buffer = append(d.buffer, line)
		case continuingArr:
			d.continuing = nil
			d.parseArrayElement("- "+line, d.continuingIndent)
			return
		}
	}

	ty := analyze(line)
	switch ty.(type) {
	case arrayElement:
		d.writeBuffer()
		d.parseArrayElement(line, indent)

	case singleLine:
		d.writeBuffer()
		d.parseSingleLine(line)

	case continuingLine:
		d.writeBuffer()
		d.parseContinuingLine(line, indent)

	case continuingString:
		d.continuing = continuingString{}
		d.continuingIndent = indent
		d.parseContinuingLine(parseContStr(line).(string)+":", indent)
		d.continuingRoot = d.root[len(d.root)-1][0].(*node)

	case continuingArr:
		d.continuing = continuingArr{}
		d.continuingIndent = indent
		d.continuingRoot = d.root[len(d.root)-1][0].(*node)
	}
}