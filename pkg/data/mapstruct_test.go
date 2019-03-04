package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApply(t *testing.T) {
	ExecuteAndCompare(t, "apply", &Apply{Path: NewPath("metadata", "name"), Function: prefixer})
}

func TestGet(t *testing.T) {
	get := Get{Path: NewPath("metadata", "name")}
	node, err := ReadManifestFile("../../testdata/mapstruct/get.yaml")
	assert.Nil(t, err)
	node.Accept(&get)
	assert.Equal(t, "datanode", get.ReturnValue.(*KeyNode).Value.(string))

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
