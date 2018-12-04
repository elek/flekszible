package processor

import "github.com/elek/flekszible/pkg/data"

type Namespace struct {
	DefaultProcessor
	Namespace string
}

func (processor *Namespace) BeforeResource(resource *data.Resource) {
	resource.Content.Accept(&data.Set{Path: data.NewPath("metadata", "namespace"), NewValue: processor.Namespace})
	resource.Content.Accept(&data.Set{Path: data.NewPath("subjects", ".*", "namespace"), NewValue: processor.Namespace})
}

func init() {
	prototype := Namespace{}
	ProcessorTypeRegistry.Add(&prototype)
}
