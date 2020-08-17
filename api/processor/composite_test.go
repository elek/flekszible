package processor

import (
	"github.com/elek/flekszible/api/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComposit(t *testing.T) {
	TestFromDir(t, "composite")
}

func TestCompositInherited(t *testing.T) {
	TestFromDir(t, "composite-inherited")
}

func TestCompositWithParam(t *testing.T) {
	TestFromDir(t, "composite-param")
}

func TestCompositWithAdditionalResources(t *testing.T) {
	TestFromDir(t, "composite-resources")
}

func TestParseTransformationParameters(t *testing.T) {
	a := assert.New(t)

	config := yaml.MapSlice{}
	config = config.Put("key1", "value")

	metadata := ProcessorMetadata{
		Name: "test",
		Parameters: []ProcessorParameter{
			{
				Name:     "key1",
				Required: true,
			},
		},
	}
	result, err := parseTransformationParameters(&metadata, &config)

	a.Nil(err)
	a.Equal("value", result["key1"])
}

func TestParseTransformationParametersMissingRequired(t *testing.T) {
	a := assert.New(t)

	config := yaml.MapSlice{}

	metadata := ProcessorMetadata{
		Name: "test",
		Parameters: []ProcessorParameter{
			{
				Name:     "key1",
				Required: true,
			},
		},
	}
	_, err := parseTransformationParameters(&metadata, &config)

	a.NotNil(err)
	a.Contains(err.Error(), "key1")
}

func TestParseTransformationParametersUndefinedParam(t *testing.T) {
	a := assert.New(t)

	config := yaml.MapSlice{}
	config = config.Put("key2", "test")

	metadata := ProcessorMetadata{
		Name: "test",
		Parameters: []ProcessorParameter{
			{
				Name: "key1",
			},
		},
	}
	_, err := parseTransformationParameters(&metadata, &config)

	a.NotNil(err)
	a.Contains(err.Error(), "key2")
}

func TestParseTransformationParametersDefault(t *testing.T) {
	a := assert.New(t)

	config := yaml.MapSlice{}
	config = config.Put("key2", "test")

	metadata := ProcessorMetadata{
		Name: "test",
		Parameters: []ProcessorParameter{
			{
				Name:    "key1",
				Default: "def",
			},
			{
				Name: "key2",
			},
		},
	}
	result, err := parseTransformationParameters(&metadata, &config)

	a.Nil(err)
	a.Equal("test", result["key2"])
	a.Equal("def", result["key1"])
}

func TestParseTransformationParametersDefaultAndCustom(t *testing.T) {
	a := assert.New(t)

	config := yaml.MapSlice{}
	config = config.Put("key2", "test")

	metadata := ProcessorMetadata{
		Name: "test",
		Parameters: []ProcessorParameter{
			{
				Name:    "key2",
				Default: "def",
			},
		},
	}
	result, err := parseTransformationParameters(&metadata, &config)

	a.Nil(err)
	a.Contains("test", result["key2"])
}
