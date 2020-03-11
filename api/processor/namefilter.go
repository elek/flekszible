package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
)

type NameFilter struct {
	DefaultProcessor
	Include []string
	Exclude []string
}

func (nf *NameFilter) BeforeResource(resource *data.Resource) error {
	if len(nf.Include) > 0 {
		for _, include := range nf.Include {
			if resource.Name() != include {
				resource.Metadata["exclude"] = "true"
			}
		}
	}

	if len(nf.Exclude) > 0 {
		for _, exclude := range nf.Exclude {
			if resource.Name() == exclude {
				resource.Metadata["exclude"] = "true"
			}
		}
	}
	return nil
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "NameFilter",
			Description: "Include and exclude certain resources (based on name)",
			Parameter: []ProcessorParameter{

				{
					Name:        "include",
					Description: "List of names to include. If set, all the other resources will be excluded",
					Default:     "",
				},
				{
					Name:        "exclude",
					Description: "List of resource names to include",
				},
			},
			Doc: "",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			processor := &NameFilter{}
			_, err := configureProcessorFromYamlFragment(processor, config)
			if err != nil {
				return processor, err
			}
			return processor, nil
		},
	})
}
