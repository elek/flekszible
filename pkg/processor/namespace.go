package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type Namespace struct {
	DefaultProcessor
	Namespace string
}

func (processor *Namespace) BeforeResource(resource *data.Resource) {
	resource.Content.Accept(&data.Set{Path: data.NewPath("metadata", "namespace"), NewValue: processor.Namespace})
	resource.Content.Accept(&data.Set{Path: data.NewPath("subjects", ".*", "namespace"), NewValue: processor.Namespace})
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name: "Namespace",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Namespace{}, config)
		},
	})
}
