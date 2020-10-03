package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
)

type DaemonToStatefulSet struct {
	DefaultProcessor
	Trigger Trigger
}

func (processor *DaemonToStatefulSet) Before(ctx *RenderContext, node *ResourceNode) error {

	newResources := make([]*data.Resource, 0)
	for _, resource := range node.AllResources() {
		if resource.Kind() == "DaemonSet" && processor.Trigger.active(resource) {

			serviceNode := createService(resource)

			newResources = append(newResources, &data.Resource{
				Content: serviceNode,
			})

		}
	}
	ctx.AddResources(newResources...)
	return nil
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

func (processor *DaemonToStatefulSet) BeforeResource(resource *data.Resource, location *ResourceNode) error {

	if resource.Kind() == "DaemonSet" && processor.Trigger.active(resource) {
		resource.Content.Get("kind").(*data.KeyNode).Value = "StatefulSet"

		name := resource.Name()

		spec := resource.Content.Get("spec").(*data.MapNode)
		spec.Put("serviceName", &data.KeyNode{Path: data.NewPath("spec", "serviceName"), Value: name})
		spec.Put("replicas", &data.KeyNode{Path: data.NewPath("spec", "replicas"), Value: 3})
	}
	return nil
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "DaemonToStatefulSet",
			Description: "Converts daemonset to statefulset",
			Parameters: []ProcessorParameter{
				TriggerParameter,
			},
			Doc: `Useful for minikube based environments where you may not have enough node to run a daemonset based cluster.` + TriggerDoc,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&DaemonToStatefulSet{}, config)
		},
	})
}
