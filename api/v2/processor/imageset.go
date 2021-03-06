package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
)

type Image struct {
	DefaultProcessor
	Image   string
	Trigger Trigger
}

func (processor *Image) ToString() string {
	return CreateToString("image").
		Add("image", processor.Image).
		Build()
}

func (imageSet *Image) BeforeResource(resource *data.Resource) error {
	if imageSet.Trigger.active(resource) {
		resource.Content.Accept(&data.Set{Path: data.NewPath("spec", "template", "spec", "(initC|c)ontainers", ".*", "image"), NewValue: imageSet.Image})
	}
	return nil
}

func ActivateImageSet(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Image",
			Description: "Replaces the docker image definition",
			Doc:         "Note: This transformations could also added with the `--image` CLI argument.",
			Parameters: []ProcessorParameter{
				TriggerParameter,
				{
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
