package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetName(t *testing.T) {
	node := NewMapNode(NewPath())
	metadata := node.CreateMap("metadata")
	metadata.PutValue("name", "{{ something}}-ok")

	r := Resource{Content: &node}

	assert.Equal(t, "ok", r.Name())
}

func TestReadConfigMaps(t *testing.T) {
	resources, err := ReadConfigMaps("../../testdata")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resources))

	expected, err := ReadManifestFile("../../testdata/configmaps/expected/config.yaml")
	assert.Nil(t, err)

	assert.Equal(t, expected, resources[0].Content)

	expected, err = ReadManifestFile("../../testdata/configmaps/expected/config2.yaml")
	assert.Nil(t, err)

	assert.Equal(t, expected, resources[1].Content)

}
