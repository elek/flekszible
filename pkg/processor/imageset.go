package processor

import "github.com/elek/flekszible/pkg/data"

type ImageSet struct {
	DefaultProcessor
	Image string
}

func (imageSet *ImageSet) BeforeResource(resource *data.Resource) {
	resource.Content.Accept(data.Set{Path: data.NewPath("spec", "template", "spec", "containers", "*", "image"), NewValue: imageSet.Image})
	resource.Content.Accept(data.Set{Path: data.NewPath("spec", "template", "spec", "initContainers", "*", "image"), NewValue: imageSet.Image})
}
func init() {
	prototype := ImageSet{}
	ProcessorTypeRegistry.Add(&prototype)
}