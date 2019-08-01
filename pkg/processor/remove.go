package processor

import (
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type Remove struct {
	DefaultProcessor
	Path    data.Path
	Trigger Trigger
	Yamlize bool
}

func (processor *Remove) BeforeResource(resource *data.Resource) {
	if !processor.Trigger.active(resource) {
		return
	}

	forceYaml := data.Yamlize{Path: processor.Path}

	if processor.Yamlize {
		resource.Content.Accept(&forceYaml)
	}

	target := data.SmartGetAll{Path: processor.Path.Parent()}
	resource.Content.Accept(&target)
	for _, match := range target.Result {
		switch typedTarget := match.Value.(type) {
		case *data.MapNode:

			typedTarget.Remove(processor.Path.Last())

		default:
			panic(fmt.Errorf("Unsupported value %T should point to a map element", processor.Path))
		}
	}

	if processor.Yamlize {
		forceYaml.Serialize = true
		resource.Content.Accept(&forceYaml)
	}

}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Remove",
			Description: "Remove yaml fragment from an existing k8s resources",
			Doc:         "",
			Parameter: []ProcessorParameter{
				PathParameter,
				TriggerParameter,
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Remove{}, config)
		},
	})
}
