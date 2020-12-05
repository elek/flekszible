package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
)

type PublishService struct {
	DefaultProcessor
	Trigger     Trigger
	ServiceType string         `yaml:"serviceType"`
	NodePorts   map[string]int `yaml:"nodePorts"`
}

func (processor *PublishService) ToString() string {
	return CreateToString("publishService").
		Add("serviceType", processor.ServiceType).
		Build()
}

func (processor *PublishService) Before(ctx *RenderContext, node *ResourceNode) error {
	newResources := make([]*data.Resource, 0)
	for _, resource := range node.AllResources() {
		if processor.Trigger.active(resource) && resource.Kind() == "Service" && hasNoneClusterIp(resource.Content) {
			newContent := DeepCopy(resource.Content)

			metadata := newContent.Get("metadata").(*data.MapNode)
			metadata.PutValue("name", metadata.GetStringValue("name")+"-public")
			spec := newContent.Get("spec").(*data.MapNode)
			spec.Remove("clusterIP")
			spec.PutValue("type", processor.ServiceType)
			r := data.NewResource()
			r.Content = newContent
			r.Destination = resource.Destination

			newResources = append(newResources, &r)
			if processor.ServiceType == "NodePort" && len(processor.NodePorts) > 0 {
				ports := spec.Get("ports").(*data.ListNode)
				for _, port := range ports.Children {
					portMap := port.(*data.MapNode)
					name := portMap.GetStringValue("name")
					for portName, nodePort := range processor.NodePorts {
						if portName == name {
							portMap.PutValue("nodePort", nodePort)
						}
					}
				}
			}
		}

	}
	ctx.AddResources(newResources...)
	return nil
}

func ActivatePublishService(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "PublishService",
			Description: "Creates additional service for internal services",
			Parameters: []ProcessorParameter{
				TriggerParameter,
				{
					Name:        "serviceType",
					Default:     "NodeType",
					Description: "The type of the newly created service.",
				},
				TriggerParameter,
				{
					Name:        "nodePorts",
					Description: "Key value map (string -> int) to define nodePort for the specific ports (In case of NodePort type services.)",
				},
			},
		},
		Factory: func(slice *yaml.MapSlice) (Processor, error) {
			return &PublishService{
				ServiceType: "NodePort",
			}, nil
		},
	})
}

func hasNoneClusterIp(slice *data.MapNode) bool {
	return true
}
