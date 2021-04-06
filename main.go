package main

import (
	"io/ioutil"
	"os"
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

func main() {
	txt, _ := readFile("test.yaml")

	rn := createNode("Root", 0, make([]*node, 0))

	tmp := []interface{}{rn, -1}
	root := [][]interface{}{tmp}
	d := diagnostic{root: root}

	T := new(tree)
	T.root = rn

	for _, line := range txt {
		indent, cutted := trim(line)
		//fmt.Printf("%d\n", indent)
		d.scan(cutted, indent, T)
	}

	T.travel(T.root)
}
