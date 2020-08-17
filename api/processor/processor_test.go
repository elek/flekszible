package processor

import (
	"github.com/elek/flekszible/api/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigureProcessorFromYamlFragment(t *testing.T) {
	assert := assert.New(t)

	i := Image{}
	k := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte("image: test"), &k)
	assert.Nil(err)

	_, err = configureProcessorFromYamlFragment(&i, &k)

	assert.Nil(err)
	assert.Equal("test", i.Image)
}

func TestConfigureProcessorFromYamlFragmentWrongYaml(t *testing.T) {
	assert := assert.New(t)

	i := Image{}
	k := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte("ximage: test"), &k)
	assert.Nil(err)

	_, err = configureProcessorFromYamlFragment(&i, &k)

	assert.NotNil(err)
}
