package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

type Namespace struct {
	DefaultProcessor
	Namespace string
	Force     bool

}

func (processor *Namespace) BeforeResource(resource *data.Resource) {
	pathList := []data.Path{data.NewPath("metadata", "namespace"), data.NewPath("subjects", ".*", "namespace")}
	for _, path := range pathList {
		if processor.Force {
			resource.Content.Accept(&data.Set{Path: path, NewValue: processor.Namespace})

		} else {
			resource.Content.Accept(&data.ReSet{Path: path, NewValue: processor.Namespace})

		}
	}
}

func init() {

	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Namespace",
			Description: "Use explicit namespace",
			Doc: `Note: This transformations could also added with the '--namespace' CLI argument.

Example ('transformations/set.yaml''):

'''yaml
- type: Namespace
  namespace: myns
'''
`,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			ns := &Namespace{}
			_, err := configureProcessorFromYamlFragment(ns, config)
			if err != nil {
				return ns, err
			}
			if ns.Namespace == "" {
				conf := data.CreateKubeConfig()
				currentNamespace, err := conf.ReadCurrentNamespace()
				if err != nil {
					return ns, err
				}
				ns.Namespace = currentNamespace
			}
			return ns, nil
		},
	})
}

