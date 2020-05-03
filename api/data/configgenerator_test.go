package data

import "github.com/stretchr/testify/assert"
import "testing"

func TestReadConfigMaps(t *testing.T) {
	generator := &ConfigGenerator{}
	resources, err := generator.Generate("../../testdata/configmaps", "/tmp")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resources))

	expected, err := ReadManifestFile("../../testdata/configmaps/expected/config.yaml")
	assert.Nil(t, err)

	assert.Equal(t, expected, resources[0].Content)

	expected, err = ReadManifestFile("../../testdata/configmaps/expected/config2.yaml")
	assert.Nil(t, err)

	assert.Equal(t, expected, resources[1].Content)

}
