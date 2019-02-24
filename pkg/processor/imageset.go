package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type Image struct {
	DefaultProcessor
	Image   string
	Trigger Trigger
}

func (imageSet *Image) BeforeResource(resource *data.Resource) {
	if imageSet.Trigger.active(resource) {
		resource.Content.Accept(&data.Set{Path: data.NewPath("spec", "template", "spec", "(initC|c)ontainers", ".*", "image"), NewValue: imageSet.Image})
	}
}
func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Image",
			Description: "Replaces the docker image definition",
			Doc:         "Note: This transformations could also added with the `--image` CLI argument.",
			Parameter: []ProcessorParameter{
				ProcessorParameter{
					Name:        "image",
					Description: "The docker image name to use as a replacement (eg. flokkr/hadoop:3.2.0)",
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Image{}, config)
		},
	})
}
