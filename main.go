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
	d := newYamlParser()
	for i, line := range txt {
		indent, cutted := trim(line)
		//	fmt.Printf("%d\n", indent)
		if i == 7 {
			println()
		}
		d.scan(cutted, indent)
	}

	d.tree.travel(d.tree.root)
}
