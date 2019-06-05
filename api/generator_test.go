package flekszible

import (
	"github.com/elek/flekszible/pkg/processor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiUsage(t *testing.T) {
	processors := make([]processor.Processor, 0)
	processors = append(processors, &processor.Image{Image: "test"})
	result, err := Generate("../testdata/api", processors)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))

	assert.Equal(t, "nginx-deployment", result[0].Name())
	s, err := result[0].Content.ToString()
	assert.Nil(t, err)
	assert.Contains(t, s, "image: test")
}
