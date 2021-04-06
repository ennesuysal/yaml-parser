package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

func trim(line string) (int, string) {
	i := 0
	for ; i < len(line); i+=2 {
		if line[i] != ' ' && line[i] != '\t'{
			break
		}
	}

	indent := i

	// out, _ := parseArrayElement(line)
	if isArrayElement(line) {
		indent+=1
	}

	return indent, line[i:]
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

func parseArray(arr string) []interface{} {
	rgx := `([^\]\[\s,]*)\s*[,\]]\s*`

	out, err := rgx_shortcut(rgx, arr)
	if err != nil {
		return nil
	}

	if out != nil {
		ret := make([]interface{}, 0)
		for _, x := range out {
				ret = append(ret, x[1])
		}
		return ret
	}

	return nil
}

func checkArray(p* node) *node {
	for _, x := range p.children {
		if x.ty == 1 {
			return x
		}
	}
	return nil
}

func isArrayElement(str string) bool {
	rgx := `^-\s*([^\s]*)\s*:\s*([^\s]*)\s*`
	match, _ := regexp.MatchString(rgx, str)
	return match
}

func parseArrayElement(str string)([][][]byte, error){
	rgx := `-\s*([^\s]*)\s*:\s*([^\s]*)\s*`
	out, err := rgx_shortcut(rgx, str)
	return out, err
}

func createNode(value interface{}, ty int, children []*node) *node {
	n := &node{
		value:    value,
		ty:       ty,
		children: children,
	}
	return n
}

func diagnostic(str string, root [][]interface{}, indent int, T* tree) [][]interface{} {
	single := `(.+)\s*:\s*(.+)`
	cont := `([^\s]+)\s*:\s*$`

	for lastRoot := root[len(root)-1][1].(int); indent <= lastRoot; lastRoot = root[len(root)-1][1].(int) {
			root = root[:len(root)-1]
	}

	out, _ := parseArrayElement(str)
	if out != nil {
		key := createNode(out[0][1], 0, make([]*node, 0))

		pa := checkArray(root[len(root)-1][0].(*node))

		if pa == nil {
			pa = createNode(make([]interface{}, 0), 1, nil)

			T.insert(root[len(root)-1][0].(*node), pa)
		}

		pa.value = append(pa.value.([]interface{}), key)

		if string(out[0][2]) != "" {
			child := createNode(out[0][2], 0, nil)
			T.insert(key, child)
		}

		root = append(root, []interface{}{pa, indent})
		root = append(root, []interface{}{key, indent})
		return root
	}

	out, _ = rgx_shortcut(single, str)
	if out != nil {

		key := createNode(out[0][1], 0, nil)
		value := createNode(out[0][2], 0, nil)

		T.insert(root[len(root) - 1][0].(*node), key)
		T.insert(key, value)

		return root
	}

	out, _ = rgx_shortcut(cont, str)
	if out != nil {

		key := createNode(out[0][1], 0, nil)

		T.insert(root[len(root)-1][0].(*node), key)

		root = append(root, []interface{}{key, indent})
		return root
	}

	arr := parseArray(str)
	if arr != nil {
		for i, x := range arr {
			fmt.Printf("%d.) %s\n", i, x)
		}
	}

	fmt.Println("")
	return root
}


func main() {
	txt, _ := readFile("test.yaml")

	rn := createNode("Root", 0, make([]*node, 0))

	tmp := []interface{}{rn, -1}
	root := [][]interface{}{tmp}

	T := new(tree)
	T.root = rn

	for _, line := range txt {
		indent, cutted := trim(line)
		//fmt.Printf("%d\n", indent)
		root = diagnostic(cutted, root, indent, T)
	}

	T.travel(T.root)
}