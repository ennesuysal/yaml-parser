// Yaml Parser Sample Usage
// Author: Enes Uysal
package main

import "fmt"

func main() {
	yp := newYamlParser("test.yaml")
	n := yp.tree.getNodeValue("jobs", "build", "steps", 0)
	fmt.Printf("%s\n", n.value)
	yp.tree.setNodeValue("enes", "jobs", "build", "steps", 0)
	n = yp.tree.getNodeValue("jobs", "build", "steps", 0)
	fmt.Printf("%s", n.value)
}
