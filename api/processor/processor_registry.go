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
	Parameters  []ProcessorParameter `yaml:"parameters"`
	Doc         string
	Resources   string //directory point to additional resources

}
type ProcessorParameter struct {
	Name        string
	Description string
	Required    bool
	Default     string
	Type        string
}

func (metadata *ProcessorMetadata) FindParam(name string) *ProcessorParameter {
	for _, param := range metadata.Parameters {
		if param.Name == name {
			return &param
		}
	}
	return nil
}

type ProcessorDefinition struct {
	Metadata ProcessorMetadata //metadata to define the name and available parameters
	Factory  ProcessorFactory  //the factory the create the struct
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
			paramDef := factory.Metadata.FindParam(key)

			if paramDef != nil && paramDef.Type == "[]string" {
				param = param.Put(key, []string{value})
			} else {
				param = param.Put(key, value)
			}
		}
		return factory.Factory(&param)
	} else {
		return nil, errors.New("No such registered transformation definition: " + name)
	}

}
