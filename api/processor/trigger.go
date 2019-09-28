package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
)

type Trigger struct {
	Definition data.Node
}

func (trigger *Trigger) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rawContent := yaml.MapSlice{}
	err := unmarshal(&rawContent)
	if err != nil {
		return err
	}
	node, err := data.ConvertToNode(rawContent, data.NewPath())
	if err != nil {
		return err
	}
	trigger.Definition = node
	return nil
}

func (trigger *Trigger) active(resource *data.Resource) bool {
	getAllKeys := data.GetKeys{}
	if trigger.Definition == nil {
		return true
	}
	trigger.Definition.Accept(&getAllKeys)
	for _, result := range getAllKeys.Result {
		getter := data.Get{Path: result.Path}
		resource.Content.Accept(&getter)
		if !getter.Found || getter.ReturnValue.(*data.KeyNode).Value != result.Value.(*data.KeyNode).Value {
			return false
		}
	}
	return true
}
