package data

import (
	"fmt"
	"github.com/elek/flekszible/pkg/yaml"
	"io/ioutil"
	"strconv"
)

func ReadString(content []byte) (*MapNode, error) {
	yamlDoc := yaml.MapSlice{}
	err := yaml.Unmarshal(content, &yamlDoc)
	if err != nil {
		return nil, err
	}

	node, err := ConvertToNode(yamlDoc, NewPath())
	if err != nil {
		return nil, err
	}
	result := node.(*MapNode)
	return result, nil
}

func ReadFile(file string) (*MapNode, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ReadString(data)
}

func ConvertToNode(object interface{}, path Path) (Node, error) {
	switch object := object.(type) {
	case yaml.MapSlice:
		result := NewMapNode(path)
		for _, pair := range object {
			key := pair.Key
			value := pair.Value
			if key == nil {
				return nil, fmt.Errorf("Key is nil at %s", path.ToString())
			}
			node, err := ConvertToNode(value, path.Extend(key.(string)))
			if err != nil {
				return nil, err
			}
			result.Put(key.(string), node)
		}

		return &result, nil

	case []interface{}:
		result := NewListNode(path)
		for ix, value := range object {
			pathSegment := strconv.Itoa(ix)
			if valueMap, ok := value.(yaml.MapSlice); ok {
				if name, ok := valueMap.Get("name"); ok {
					pathSegment = name.(string)
				}
			}
			node, err := ConvertToNode(value, path.Extend(pathSegment))
			if err != nil {
				return nil, err
			}
			result.Append(node)
		}
		return &result, nil
	case string:
		res := NewKeyNode(path, object)
		return &res, nil
	case bool:
		res := NewKeyNode(path, object)
		return &res, nil
	case int:
		res := NewKeyNode(path, object)
		return &res, nil
	default:
		return nil, fmt.Errorf("I don't know about type %T! %s\n", object, path.ToString())
	}
}
