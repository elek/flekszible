package processor

import (
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type Add struct {
	DefaultProcessor
	Path    data.Path
	Trigger Trigger
	Value   interface{}
}

func (processor *Add) BeforeResource(resource *data.Resource) {
	if !processor.Trigger.active(resource) {
		return
	}
	switch typedValue := processor.Value.(type) {
	case yaml.MapSlice:
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
			default:
				panic(fmt.Errorf("Unsupported value adding %T to %T", match.Value, processor.Value))
			}
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
				panic(fmt.Errorf("Unsupported value adding %T to %T %s", match.Value, processor.Value, resource.Filename))
			}
		}
	default:
		panic(fmt.Errorf("Unsupported value adding %T", processor.Value))
	}
}

func init() {
	prototype := Add{}
	ProcessorTypeRegistry.Add(&prototype)
}
