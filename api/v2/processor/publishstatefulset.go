package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
)

func ActivatePublishStatefulset(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "PublishStatefulSet",
			Description: "Creates additional NodeType service for StatefulSet internal services",
			Parameters: []ProcessorParameter{
				TriggerParameter,
				{
					Name:        "ports",
					Description: "Key value map (string -> int) to define nodePort for the specific ports.",
				},
			},
		},
		Factory: func(slice *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&PublishService{NodePorts: make(map[string]int), ServiceType: "NodePort"}, slice)
		},
	})
}

func DeepCopy(src *data.MapNode) *data.MapNode {
	content, err := src.ToString()
	if err != nil {
		panic(err)
	}
	mapNode, err := data.ReadManifestString([]byte(content))
	if err != nil {
		panic(err)
	}
	return mapNode
}
