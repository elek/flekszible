package data

import "github.com/stretchr/testify/assert"
import "testing"

func TestReadConfigMaps(t *testing.T) {
	generator := &ConfigGenerator{}
	resources, err := generator.Generate("../../testdata/configmaps", "/tmp")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resources))

	nameConfig, err := ReadManifestFile("../../testdata/configmaps/expected/config.yaml")
	assert.Nil(t, err)

	otherConfig, err := ReadManifestFile("../../testdata/configmaps/expected/config2.yaml")
	assert.Nil(t, err)

	if resources[0].Name() == "name" {
		assert.Equal(t, nameConfig, resources[0].Content)
		assert.Equal(t, otherConfig, resources[1].Content)
	} else {
		assert.Equal(t, nameConfig, resources[1].Content)
		assert.Equal(t, otherConfig, resources[0].Content)
	}

}
