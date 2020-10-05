package processor

import (
	"testing"

	"github.com/elek/flekszible/api/v2/data"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromFile(t *testing.T) {
	registry := NewRegistry()
	ActivateAdd(registry)

	testfile := "../../testdata/processors/prometheus.yaml"
	processors, err := registry.ReadProcessorDefinitionFile(testfile)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(processors))
	first := processors[0]
	processorAdd, ok := first.(*Add)
	assert.True(t, ok, "The first processor was not an Add")
	assert.EqualValues(t, data.NewPath("spec", "templalte", "metadata", "annotations"), processorAdd.Path)

}

func TestSimpleCreate(t *testing.T) {
	registry := NewRegistry()
	ActivateImageSet(registry)
	parameters := map[string]string{"image": "test"}
	proc, err := registry.Create("image", parameters)
	assert.Nil(t, err)
	assert.NotNil(t, proc)
	image := proc.(*Image)
	assert.Equal(t, "test", image.Image)
}
