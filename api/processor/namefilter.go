package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"strings"
)

type NameFilter struct {
	DefaultProcessor
	Include []string
	Exclude []string
}

func (nf *NameFilter) ToString() string {
	return CreateToString("namefilter").
		Add("include", strings.Join(nf.Include, ",")).
		Add("exclude", strings.Join(nf.Exclude, ",")).
		Build()
}

func (nf *NameFilter) BeforeResource(resource *data.Resource) error {
	exclude := true
	if len(nf.Include) > 0 {
		for _, include := range nf.Include {
			if resource.Name() == include {
				exclude = false
			}
		}
	} else {
		exclude = false
	}

	if len(nf.Exclude) > 0 {
		for _, excludeRule := range nf.Exclude {
			if resource.Name() == excludeRule {
				exclude = true
			}
		}
	}
	if exclude {
		resource.Metadata["exclude"] = "true"
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
					Type:        "[]string",
				},
				{
					Name:        "exclude",
					Description: "List of resource names to include",
					Type:        "[]string",
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
