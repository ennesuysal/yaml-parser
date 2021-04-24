## Ordered Yaml Parser
### Usage:

#### Download Package
    go get github.com/ennesuysal/yaml-parser

#### Import Package to Your Project
    import "yaml-parser"

#### Create an instance
    yp := yamlParser.NewYamlParser("test.yaml")

#### Get Value
    n := yp.Tree.GetNodeValue("build", "docker", 0)
    value := n.Value

#### Set Value
    yp.Tree.SetNodeValue("newValue", "build", "docker", 0)

#### Write to file
    yp.Tree.WriteFile("output.yaml")