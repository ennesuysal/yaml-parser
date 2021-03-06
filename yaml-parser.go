// Author: Enes Uysal
package yamlParser

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

func NewYamlParser(filePath string) *Diagnostic {
	txt, _ := readFile(filePath)
	d := newYamlHelper()
	for _, line := range txt {
		indent, trimmed := trim(line)
		d.scan(trimmed, indent)
	}
	return d
}
