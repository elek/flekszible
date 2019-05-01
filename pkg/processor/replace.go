package processor

import (
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type Replace struct {
	DefaultProcessor
	Path    data.Path
	Trigger Trigger
	Value   interface{}
}

func (processor *Replace) BeforeResource(resource *data.Resource) {
	if !processor.Trigger.active(resource) {
		return
	}

	target := data.SmartGetAll{Path: processor.Path.Parent()}
	resource.Content.Accept(&target)
	for _, match := range target.Result {
		switch typedTarget := match.Value.(type) {
		case *data.MapNode:
			node, err := data.ConvertToNode(processor.Value, processor.Path)
			if err != nil {
				panic(err)
			}
			typedTarget.Put(processor.Path.Last(), node)
		default:
			panic(fmt.Errorf("Unsupported value adding %T to %T", processor.Value, match.Value))
		}
	}
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Replace",
			Description: "Replace a yaml subtree with an other one.",
			Doc:         addDocReplace,
			Parameter: []ProcessorParameter{
				PathParameter,
				TriggerParameter,
				{
					Name:        "value",
					Description: "A yaml struct to replace the defined value",
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Replace{}, config)
		},
	})
}

var addDocReplace = `#### Supported value types

| Type of the destination node (selected by 'Path') | Type of the 'Value' | Supported
|---------------------------------------------------|---------------------|------------
| MapElement                                        | Any Yaml            | Yes
'''` + TriggerDoc + PathDoc
