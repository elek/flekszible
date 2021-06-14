package processor

import (
	"fmt"

	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/pkg/errors"
)

type Remove struct {
	DefaultProcessor
	Path    data.Path
	Trigger Trigger
	Yamlize bool
}

func (processor *Remove) BeforeResource(resource *data.Resource) error {
	if !processor.Trigger.active(resource) {
		return nil
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
		case *data.ListNode:
			indexToRemove := -1
			for ix, child := range typedTarget.Children {
				switch typedChild := child.(type) {
				case *data.MapNode:
					if typedChild.GetStringValue("name") == processor.Path.Last() {
						indexToRemove = ix
					}
				}
			}
			if indexToRemove >= 0 {
				copy(typedTarget.Children[indexToRemove:], typedTarget.Children[indexToRemove+1:])
				typedTarget.Children[len(typedTarget.Children)-1] = nil
				typedTarget.Children = typedTarget.Children[:len(typedTarget.Children)-1]
			} else {
				return errors.New(fmt.Sprintf("Path %s points a list to remove an element, but no map element with name=%s", processor.Path.ToString(), processor.Path.Last()))
			}

		default:
			return errors.New(fmt.Sprintf("Unsupported path to remove: %s should point to a map type element but %T", processor.Path.ToString(), typedTarget))
		}
	}

	if processor.Yamlize {
		forceYaml.Serialize = true
		resource.Content.Accept(&forceYaml)
	}
	return nil
}

func ActivateRemove(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Remove",
			Description: "Remove yaml fragment from an existing k8s resources",
			Doc:         "",
			Parameters: []ProcessorParameter{
				PathParameter,
				TriggerParameter,
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Remove{}, config)
		},
	})
}
