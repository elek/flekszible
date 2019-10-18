package processor

import (
	"testing"

	"github.com/elek/flekszible/api/data"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromFile(t *testing.T) {
	testfile := "../../testdata/processors/prometheus.yaml"
	processors, err := ReadProcessorDefinitionFile(testfile)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(processors))
	first := processors[0]
	processorAdd, ok := first.(*Add)
	assert.True(t, ok, "The first processor was not an Add")
	assert.EqualValues(t, data.NewPath("spec", "templalte", "metadata", "annotations"), processorAdd.Path)

}

func TestSimpleCreate(t *testing.T) {
	parameters := map[string]string{"image": "test"}
	proc, err := ProcessorTypeRegistry.Create("image", parameters)
	assert.Nil(t, err)
	assert.NotNil(t, proc)
	image := proc.(*Image)
	assert.Equal(t, "test", image.Image)
}
