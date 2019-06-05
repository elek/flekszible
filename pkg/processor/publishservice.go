package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type PublishService struct {
	DefaultProcessor
	Trigger Trigger
	Type    string
}

func (processor *PublishService) Before(ctx *RenderContext, resources []*data.Resource) {
	newResources := make([]*data.Resource, 0)
	for _, resource := range resources {
		if processor.Trigger.active(resource) && resource.Kind() == "Service" && hasNoneClusterIp(resource.Content) {
			newContent := DeepCopy(resource.Content)

			metadata := newContent.Get("metadata").(*data.MapNode)
			metadata.PutValue("name", metadata.GetStringValue("name")+"-public")
			spec := newContent.Get("spec").(*data.MapNode)
			spec.Remove("clusterIP")
			spec.PutValue("type", processor.Type)
			r := data.Resource{
				Content:     newContent,
				Destination: resource.Destination,
			}
			newResources = append(newResources, &r)
		}

	}
	ctx.AddResources(newResources...)
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "PublishService",
			Description: "Creates additional service for internal services",
			Parameter: []ProcessorParameter{
				TriggerParameter,
				{
					Name:        "type",
					Default:     "NodeType",
					Description: "The type of the newly created service.",
				},
			},
		},
		Factory: func(slice *yaml.MapSlice) (Processor, error) {
			return &PublishService{
				Type: "NodePort",
			}, nil
		},
	})
}
