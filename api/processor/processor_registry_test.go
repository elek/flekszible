package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/stretchr/testify/assert"
	"testing"
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
