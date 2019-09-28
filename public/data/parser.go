package data

import (
	"fmt"
	"github.com/elek/flekszible/public/yaml"
	"io/ioutil"
	"strconv"
)

func ReadManifestString(content []byte) (*MapNode, error) {
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

func ReadManifestFile(file string) (*MapNode, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ReadManifestString(data)
}

//Converts internal node tree to Yaml
func ConvertToYaml(root Node) interface{} {
	switch object := root.(type) {
	case *MapNode:
		slice := yaml.MapSlice{}
		for _, key := range object.keys {
			convertedValue := ConvertToYaml(object.Get(key))
			slice = slice.Put(key, convertedValue)
		}
		return slice
	case *KeyNode:
		return object.Value
	case *ListNode:
		list := make([]interface{}, 0)
		for _, item := range object.Children {
			list = append(list, ConvertToYaml(item))
		}
		return list
	}
	return nil
}

//Converts raw Yaml structure to file to internal node tree
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
	case nil:
		res := NewKeyNode(path, nil)
		return &res, nil
	default:
		return nil, fmt.Errorf("I don't know about type %T! %s\n", object, path.ToString())
	}
}
