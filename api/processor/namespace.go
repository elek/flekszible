package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
)

type Namespace struct {
	DefaultProcessor
	Namespace          string
	Force              bool
	ClusterRoleSupport bool
}

func (processor *Namespace) ToString() string {
	return CreateToString("namespace").
		Add("namespace", processor.Namespace).
		AddBool("cluserRoleSupport", processor.ClusterRoleSupport).
		AddBool("force", processor.Force).
		Build()
}

func (processor *Namespace) BeforeResource(resource *data.Resource) error {
	pathList := []data.Path{data.NewPath("metadata", "namespace"), data.NewPath("subjects", ".*", "namespace")}
	for _, path := range pathList {
		if processor.Force {
			resource.Content.Accept(&data.Set{Path: path, NewValue: processor.Namespace})

		} else {
			resource.Content.Accept(&data.ReSet{Path: path, NewValue: processor.Namespace})

		}
	}
	if processor.ClusterRoleSupport {
		if resource.Kind() == "ClusterRole" {
			namePath := data.NewPath("metadata", "name")
			name := resource.Name()
			resource.DestinationFileName = CreateOutputFileName(name, "ClusterRole")
			resource.Content.Accept(&data.Set{Path: namePath, NewValue: name + "-" + processor.Namespace})
		}
		if resource.Kind() == "ClusterRoleBinding" {
			namePath := data.NewPath("metadata", "name")
			name := resource.Name()
			resource.DestinationFileName = CreateOutputFileName(name, "ClusterRoleBinding")
			resource.Content.Accept(&data.Set{Path: namePath, NewValue: name + "-" + processor.Namespace})
		}
		if resource.Kind() == "ClusterRoleBinding" {
			namePath := data.NewPath("roleRef", "name")
			get := &data.Get{Path: namePath}
			resource.Content.Accept(get)
			resource.Content.Accept(&data.Set{Path: namePath, NewValue: get.ValueAsString() + "-" + processor.Namespace})
		}
	}
	return nil
}

func init() {

	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Namespace",
			Description: "Use explicit namespace",
			Parameter: []ProcessorParameter{

				{
					Name:        "namespace",
					Description: "The namespace to use in the k8s resources. If empty, the current namespace will be used (from ~/.kube/config or $KUBECONFIG)",
					Default:     "",
				},
				{
					Name:        "force",
					Description: "If false (default) only the existing namespace attributes will be changed. If yes, namespace will be added to all the resources.",
					Default:     "false",
				},
				{
					Name:        "clusterrolesupport",
					Description: "If true, the created cluster roles and cluster role bindings will be postfixed by the namespace to guarantee multitenancy.",
					Default:     "true",
				},
			},
			Doc: `Note: This transformations could also added with the '--namespace' CLI argument.

Example):

'''yaml
- type: Namespace
  namespace: myns
'''
`,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			ns := &Namespace{}
			ns.ClusterRoleSupport = true
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
