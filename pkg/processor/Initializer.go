package processor

import (
	"github.com/elek/flekszible/pkg/data"
)

type Initializer struct {
	DefaultProcessor
}

func (processor *Initializer) BeforeResource(resource *data.Resource) {
	processor.ensureList(resource, data.NewPath("spec", "template", "spec", ".*ontainers", ".*"), "volumeMounts")
	processor.ensureList(resource, data.NewPath("spec", "template", "spec", ".*ontainers", ".*"), "env")
	processor.ensureList(resource, data.NewPath("spec", "template", "spec", ".*ontainers", ".*"), "envFrom")

	processor.ensureList(resource, data.NewPath("spec", "template", "spec"), "volumes")
	processor.ensureMap(resource, data.NewPath("metadata"), "labels")
	processor.ensureMap(resource, data.NewPath("metadata"), "annotations")
	processor.ensureMap(resource, data.NewPath("spec","template","metadata"), "annotations")
	processor.ensureMap(resource, data.NewPath("spec","template","metadata"), "labels")
}

func (processor *Initializer) ensureList(resource *data.Resource, path data.Path, name string) {
	get := data.GetAll{Path: path}
	resource.Content.Accept(&get)
	for ix, _ := range get.Result {
		item := get.Result[ix]
		mn := item.Value.(*data.MapNode)
		if !mn.HasKey(name) {
			newNodeItem := data.NewListNode(item.Path.Extend(name))
			mn.Put(name, &newNodeItem)
		}
	}
}

func (processor *Initializer) ensureMap(resource *data.Resource, path data.Path, name string) {
	get := data.GetAll{Path: path}
	resource.Content.Accept(&get)
	for ix, _ := range get.Result {
		item := get.Result[ix]
		mn := item.Value.(*data.MapNode)
		if !mn.HasKey(name) {
			newNodeItem := data.NewMapNode(item.Path.Extend(name))
			mn.Put(name, &newNodeItem)
		}
	}
}

//func (*Initializer) BeforeMap(path data.Path, object yaml.MapSlice) yaml.MapSlice {
//	if path.MatchSegments("metadata") {
//		object = ensureMap(object, "labels")
//		object = ensureMap(object, "annotations")
//	}
//	if path.MatchSegments("spec", "template", "spec", "containers", "*") {
//		object = ensureArray(object, "env")
//		object = ensureArray(object, "volumeMounts")
//	}
//	if path.MatchSegments("spec", "template", "spec", "initContainers", "*") {
//		object = ensureArray(object, "env")
//		object = ensureArray(object, "volumeMounts")
//	}
//	if path.MatchSegments("spec", "template", "spec") {
//		object = ensureArray(object, "volumes")
//	}
//	if path.MatchSegments("spec", "template", "spec", "containers", "*") ||
//		path.MatchSegments("spec", "template", "spec", "initContainers", "*") {
//		object = ensureArray(object, "envFrom")
//	}
//	return object
//}
//
//func ensureArray(dict yaml.MapSlice, key string) yaml.MapSlice {
//	if _, ok := dict.Get(key); !ok {
//		return dict.Put(key, make([]interface{}, 0))
//	}
//	return dict
//}
//
//func ensureMap(dict yaml.MapSlice, key string) yaml.MapSlice {
//	if _, ok := dict.Get(key); !ok {
//		return dict.Put(key, yaml.MapSlice{})
//	}
//	return dict
//}

func init() {
	prototype := Initializer{}
	ProcessorTypeRegistry.Add(&prototype)
}