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

	node := ConvertToNode(yamlDoc, NewPath())
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

func ConvertToNode(object interface{}, path Path) Node {
	switch object := object.(type) {
	case yaml.MapSlice:
		result := NewMapNode(path)
		for _, pair := range object {
			key := pair.Key
			value := pair.Value
			result.Put(key.(string), ConvertToNode(value, path.Extend(key.(string))))
		}


		return &result

	case []interface{}:
		result := NewListNode(path)
		for ix, value := range (object) {
			pathSegment := strconv.Itoa(ix)
			if valueMap, ok := value.(yaml.MapSlice); ok {
				if name, ok := valueMap.Get("name"); ok {
					pathSegment = name.(string)
				}
			}
			result.Append(ConvertToNode(value, path.Extend(pathSegment)))
		}
		return &result
	case string:
		res := NewKeyNode(path, object)
		return &res
	case bool:
		res := NewKeyNode(path, object)
		return &res
	case int:
		res := NewKeyNode(path, object)
		return &res
	default:
		panic(fmt.Errorf("I don't know about type %T! %s\n", object, path.ToString()))
	}
}
