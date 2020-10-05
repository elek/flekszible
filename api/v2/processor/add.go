package processor

import (
	"errors"
	"fmt"
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"strconv"
)

type Add struct {
	DefaultProcessor
	Path    data.Path
	Trigger Trigger
	Value   interface{}
	Yamlize bool
}

func (add *Add) ToString() string {

	return CreateToString("add").
		Add("path", add.Path.ToString()).
		AddValue("value", add.Value).
		Build()
}

func (processor *Add) BeforeResource(resource *data.Resource) error {
	if !processor.Trigger.active(resource) {
		return nil
	}
	switch typedValue := processor.Value.(type) {
	case yaml.MapSlice:
		forceYaml := data.Yamlize{Path: processor.Path}

		if processor.Yamlize {
			resource.Content.Accept(&forceYaml)
		}

		target := data.SmartGetAll{Path: processor.Path}
		resource.Content.Accept(&target)
		for _, match := range target.Result {
			switch typedTarget := match.Value.(type) {
			case *data.MapNode:
				node, err := data.ConvertToNode(typedValue, match.Path)
				if err != nil {
					panic(err)
				}
				mapNode := node.(*data.MapNode)
				for _, key := range mapNode.Keys() {
					typedTarget.Put(key, mapNode.Get(key))
				}
			case *data.ListNode:
				node, err := data.ConvertToNode(typedValue, match.Path.Extend(strconv.Itoa(typedTarget.Len())))
				if err != nil {
					panic(err)
				}
				typedTarget.Append(node)

			default:
				panic(fmt.Errorf("Unsupported value adding %T to %T", processor.Value, match.Value))
			}
		}

		if processor.Yamlize {
			forceYaml.Serialize = true
			resource.Content.Accept(&forceYaml)
		}
	case []interface{}:
		target := data.SmartGetAll{Path: processor.Path}
		resource.Content.Accept(&target)
		for _, match := range target.Result {
			switch typedTarget := match.Value.(type) {
			case *data.ListNode:
				node, err := data.ConvertToNode(typedValue, match.Path)
				if err != nil {
					panic(err)
				}
				nodeList := node.(*data.ListNode)
				for _, childNode := range nodeList.Children {
					typedTarget.Append(childNode)
				}
			default:
				return errors.New(fmt.Sprintf("Unsupported value adding %T to %T (%s) %s", processor.Value, match.Value, target.Path.ToString(), resource.Filename))
			}
		}
	default:
		return errors.New(fmt.Sprintf("Unsupported value adding %T", processor.Value))
	}
	resource.Content.Accept(&data.FixPath{CurrentPath: data.RootPath()})
	return nil
}

func ActivateAdd(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Add",
			Description: "Extends yaml fragment to an existing k8s resources",
			Doc:         addDoc,
			Parameters: []ProcessorParameter{
				PathParameter,
				TriggerParameter,
				{
					Name:        "value",
					Description: "A yaml struct to add to the defined path",
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Add{}, config)
		},
	})
}

var addDoc = `#### Supported value types

| Type of the destination node (selected by 'Path') | Type of the 'Value' | Supported
|---------------------------------------------------|---------------------|------------
| Map                                               | Map                 | Yes
| Array                                             | Array               | Yes
| Array                                             | Map                 | Yes

#### Example

'''
- type: Add
  path:
  - metadata
  - annotations
  value:
     flekszible: generated
'''` + TriggerDoc + PathDoc
