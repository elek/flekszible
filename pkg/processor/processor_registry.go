package processor

import (
	"github.com/elek/flekszible/pkg/yaml"
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
	pt.TypeMap[definition.Metadata.Name] = definition
}


