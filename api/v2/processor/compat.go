package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
)

type Compat struct {
	DefaultProcessor
	Version string
}

func (processor *Compat) BeforeResource(resource *data.Resource) error {
	if processor.Version == "1.17" {
		if resource.Kind() == "Deployment" && resource.Get(data.NewPath("apiVersion")) == "extensions/v1beta1" {
			set := data.Set{Path: data.NewPath("apiVersion"), NewValue: "apps/v1"}
			resource.Content.Accept(&set)

			spec := resource.Content.Get("spec").(*data.MapNode)
			if !spec.HasKey("selector") {
				get := data.Get{Path: data.NewPath("spec", "template", "metadata", "labels")}
				resource.Content.Accept(&get)
				if get.Found {
					selector := spec.CreateMap("selector")
					matchLabels := selector.CreateMap("matchLabels")
					for key, value := range get.ReturnValue.(*data.MapNode).ToMap() {
						matchLabels.PutValue(key, value)
					}
				}
			}
		}
	}
	return nil
}

func init() {

	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Compat",
			Description: "Kubernetes compatibilty converted",
			Parameters: []ProcessorParameter{

				{
					Name:        "version",
					Description: "Version of the Kubernetes",
					Default:     "",
				},
			},
			Doc: `This transformation tries to handle API changes between kubernetes version`,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			compat := &Compat{}
			_, err := configureProcessorFromYamlFragment(compat, config)
			return compat, err
		},
	})
}
