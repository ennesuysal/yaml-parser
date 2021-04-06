package main

import "regexp"

func (d *diagnostic) parseContinuingLine(line string, indent int) {
	out, _ := rgxShortcut(continuingLineRgx, line)
	if out != nil {
		key := createNode(out[0][1], 0, nil)
		d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		d.root = append(d.root, []interface{}{key, indent})
	}
}

func (d *diagnostic) parseSingleLine(line string) {
	out, _ := rgxShortcut(singleLineRgx, line)
	if out != nil {

		key := createNode(out[0][1], 0, nil)
		value := createNode(out[0][2], 0, nil)

		d.tree.insert(d.root[len(d.root)-1][0].(*node), key)
		d.tree.insert(key, value)
	}
}

func (d *diagnostic) parseArrayElement(line string, indent int) {
	out, _ := rgxShortcut(arrayElementRgx, line)
	if out != nil {
		key := createNode(out[0][1], 0, make([]*node, 0))

		pa := checkArray(d.root[len(d.root)-1][0].(*node))

		if pa == nil {
			pa = createNode(make([]interface{}, 0), 1, nil)

			d.tree.insert(d.root[len(d.root)-1][0].(*node), pa)
		}

		pa.value = append(pa.value.([]interface{}), key)

		if string(out[0][2]) != "" {
			child := createNode(out[0][2], 0, nil)
			d.tree.insert(key, child)
		}

		d.root = append(d.root, []interface{}{pa, indent})
		d.root = append(d.root, []interface{}{key, indent})
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
