package processor

import (
	"errors"
	"strings"

	"github.com/elek/flekszible/api/yaml"
)

var ProcessorTypeRegistry ProcessorTypes

type ProcessorMetadata struct {
	Name        string
	Description string
	Parameter   []ProcessorParameter
	Doc         string
}
type ProcessorParameter struct {
	Name        string
	Description string
	Required    bool
	Default     string
}

type ProcessorDefinition struct {
	Metadata ProcessorMetadata //metadata to define the name and available parameters
	Factory  ProcessorFactory  //the factory the create the stryct
}

type ProcessorFactory = func(slice *yaml.MapSlice) (Processor, error)

//The main processor registry
type ProcessorTypes struct {
	TypeMap map[string]ProcessorDefinition
}

func (pt *ProcessorTypes) Add(definition ProcessorDefinition) {
	if pt.TypeMap == nil {
		pt.TypeMap = make(map[string]ProcessorDefinition)
	}
	pt.TypeMap[strings.ToLower(definition.Metadata.Name)] = definition
}

func (pt *ProcessorTypes) Create(name string, parameters map[string]string) (Processor, error) {
	if factory, ok := pt.TypeMap[name]; ok {
		param := yaml.MapSlice{}
		for key, value := range parameters {
			param = param.Put(key, value)
		}
		return factory.Factory(&param)
	} else {
		return nil, errors.New("No such registered transformation definition: " + name)
	}

}
