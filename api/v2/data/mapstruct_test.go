package data

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	ExecuteAndCompare(t, "apply", &Apply{Path: NewPath("metadata", "name"), Function: prefixer})
}

func TestYamlize(t *testing.T) {
	path := NewPath("data", "dashboards.yaml", "providers", "default", "type")
	get := Get{Path: path}

	node, err := ReadManifestFile("../../testdata/mapstruct/yamlize.yaml")
	assert.Nil(t, err)
	yamlize := &Yamlize{Path: path}
	node.Accept(yamlize)
	node.Accept(&get)

	assert.True(t, get.Found)
	yamlize.Serialize = true
	node.Accept(yamlize)
	keyNode := node.Get("data").(*MapNode).Get("dashboards.yaml").(*KeyNode)
	strings.Contains("path: /etc/dashboards", keyNode.Value.(string))
}
func TestGet(t *testing.T) {
	get := Get{Path: NewPath("metadata", "name")}
	node, err := ReadManifestFile("../../testdata/mapstruct/get.yaml")
	assert.Nil(t, err)
	node.Accept(&get)
	assert.Equal(t, "datanode", get.ReturnValue.(*KeyNode).Value.(string))

}

func TestReSetReal(t *testing.T) {
	path := NewPath("metadata", "name")
	value := "n1"
	reset1 := ReSet{Path: path, NewValue: value}
	node, err := ReadManifestFile("../../testdata/mapstruct/reset.yaml")
	assert.Nil(t, err)
	node.Accept(&reset1)

	get1 := Get{Path: path}
	node.Accept(&get1)

	assert.True(t, get1.Found)
	assert.Equal(t, value, get1.ReturnValue.(*KeyNode).Value)
}

func TestReSetMissing(t *testing.T) {
	path := NewPath("metadata", "namespace")
	reset1 := ReSet{Path: path, NewValue: "n1"}
	node, err := ReadManifestFile("../../testdata/mapstruct/reset.yaml")
	assert.Nil(t, err)
	node.Accept(&reset1)

	get1 := Get{Path: path}
	node.Accept(&get1)

	assert.False(t, get1.Found)
}

func TestSmartGet(t *testing.T) {
	get := SmartGetAll{Path: NewPath("metadata", "annotations")}
	root := NewMapNode(NewPath())
	root.PutValue("test", "value")

	root.Accept(&get)

	assert.Equal(t, 1, len(get.Result))
	expected := NewMapNode(NewPath("metadata", "annotations"))
	assert.Equal(t, &expected, get.Result[0].Value)

}

func ExecuteAndCompare(t *testing.T, name string, visitor Visitor) {
	node, err := ReadManifestFile("../../testdata/mapstruct/" + name + ".yaml")
	assert.Nil(t, err)

	expected, err := ReadManifestFile("../../testdata/mapstruct/" + name + "_expected.yaml")
	assert.Nil(t, err)

	node.Accept(visitor)
	assert.Equal(t, expected, node)

	node.Accept(PrintVisitor{})
}

func prefixer(value interface{}) interface{} {
	return "xxx-" + value.(string)
}

func TestFixPath(t *testing.T) {
	n := NewMapNode(NewPath())
	childMap := n.CreateMap("child1")
	childMap.PutValue("key", "value")

	childList := n.CreateList("list")
	childList.AddValue("asd")
	childList.AddValue("bsd")
	maplist := n.CreateList("maplist")
	m1 := maplist.CreateMap()
	m1.PutValue("name", "n1")
	m1.PutValue("k1", "v2")
	fp := FixPath{}
	n.Accept(PrintVisitor{})
	println("--------------")
	n.Accept(&fp)
	n.Accept(PrintVisitor{})
}

func TestNodeFromPathValue(t *testing.T) {
	path := NewPath("metadata", "interface")
	node := NodeFromPathValue(path, "something")

	g := Get{Path: path}
	node.Accept(&g)

	assert.True(t, g.Found)
	assert.Equal(t, "something", g.ValueAsString())
}

func TestToMap(t *testing.T) {
	n := NewMapNode(NewPath())
	childMap := n.CreateMap("child1")
	childMap.PutValue("key", "value")

	toMap := n.ToMap()

	value := toMap["child1"].(map[string]interface{})["key"]
	assert.Equal(t, "value", value)
}
