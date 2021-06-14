package data

import (
	"io/ioutil"
	"testing"

	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	node, err := ReadManifestFile("../../testdata/parser/datanode.yaml")
	assert.Nil(t, err)
	node.Accept(PrintVisitor{})
}

func TestConvertToYaml(t *testing.T) {
	yamlBytes, err := ioutil.ReadFile("../../testdata/yaml/ss.yaml")
	assert.Nil(t, err)
	yamlDoc := yaml.MapSlice{}
	err = yaml.Unmarshal(yamlBytes, &yamlDoc)
	assert.Nil(t, err)
	node, err := ConvertToNode(yamlDoc, NewPath())
	assert.Nil(t, err)
	result := ConvertToYaml(node)
	assert.Equal(t, yamlDoc, result)
}

func TestConvertToYamlWithNull(t *testing.T) {
	yamlBytes, err := ioutil.ReadFile("../../testdata/yaml/ss-with-null.yaml")
	assert.Nil(t, err)
	yamlDoc := yaml.MapSlice{}
	err = yaml.Unmarshal(yamlBytes, &yamlDoc)
	assert.Nil(t, err)
	node, err := ConvertToNode(yamlDoc, NewPath())
	assert.Nil(t, err)
	result := ConvertToYaml(node)
	assert.Equal(t, yamlDoc, result)
}
