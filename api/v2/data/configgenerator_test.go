package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigMaps(t *testing.T) {
	generator := &ConfigImporter{}
	resources, err := generator.Generate("../../testdata/configmaps", "/tmp")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resources))

	otherConfig, err := ReadManifestFile("../../testdata/configmaps/expected/config2.yaml")
	assert.Nil(t, err)

	if resources[0].Name() == "name" {
		assert.Equal(t, otherConfig, resources[1].Content)
	} else {
		assert.Equal(t, otherConfig, resources[0].Content)
	}

}
