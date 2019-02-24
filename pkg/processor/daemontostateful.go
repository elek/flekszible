package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type DaemonToStatefulSet struct {
	DefaultProcessor
	Trigger Trigger
}

func (processor *DaemonToStatefulSet) Before(ctx *RenderContext, resources []data.Resource) {

	newResources := make([]data.Resource, 0)
	for _, resource := range resources {
		if resource.Kind() == "DaemonSet" && processor.Trigger.active(&resource) {

			serviceNode := createService(&resource)

			newResources = append(newResources, data.Resource{
				Content: serviceNode,
			})

		}
	}
	ctx.AddResources(newResources...)
}
func createService(resource *data.Resource) *data.MapNode {
	root := data.NewMapNode(data.NewPath())
	root.PutValue("apiVersion", "v1")
	root.PutValue("kind", "Service")
	metadata := root.CreateMap("metadata")
	metadata.PutValue("name", resource.Name())

	spec := root.CreateMap("spec")
	spec.PutValue("clusterIP", "None")

	selector := data.Get{Path: data.NewPath("spec", "selector", "matchLabels")}
	resource.Content.Accept(&selector)
	if selector.Found {
		newSelector := spec.CreateMap("selector")
		selectorMap := selector.ReturnValue.(*data.MapNode)
		for _, key := range selectorMap.Keys() {
			value := selectorMap.Get(key).(*data.KeyNode).Value.(string)
			newSelector.PutValue(key, value)
		}
	}
	root.Accept(&data.FixPath{CurrentPath: data.NewPath()})

	return &root
}

func (processor *DaemonToStatefulSet) BeforeResource(resource *data.Resource) {

	if resource.Kind() == "DaemonSet" && processor.Trigger.active(resource) {
		resource.Content.Get("kind").(*data.KeyNode).Value = "StatefulSet"

		name := resource.Name()

		spec := resource.Content.Get("spec").(*data.MapNode)
		spec.Put("serviceName", &data.KeyNode{Path: data.NewPath("spec", "serviceName"), Value: name})
		spec.Put("replicas", &data.KeyNode{Path: data.NewPath("spec", "replicas"), Value: 3})
	}
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name: "DaemonToStatefulSet",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&DaemonToStatefulSet{}, config)
		},
	})
}
