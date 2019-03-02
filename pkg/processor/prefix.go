package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"strings"
)

type Prefix struct {
	DefaultProcessor
	Prefix    string
	HostNames []string
}

func (p *Prefix) Before(ctx *RenderContext, resources []*data.Resource) {
	p.HostNames = make([]string, 0)
	for _, resource := range resources {
		kind := resource.Kind()
		if kind == "StatefulSet" {
			name := resource.Name()

			serviceNameGet := data.Get{Path: data.NewPath("spec", "serviceName")}
			resource.Content.Accept(&serviceNameGet)
			if serviceNameGet.Found {
				hostname := name + "-0." + serviceNameGet.ValueAsString()
				p.HostNames = append(p.HostNames, hostname)
			}
		}
	}
}

func (p *Prefix) BeforeResource(resource *data.Resource) {

	content := resource.Content

	prefixName := func(original interface{}) interface{} {
		return p.Prefix + "-" + original.(string)
	}

	prefixHostName := func(original interface{}) interface{} {
		result := original.(string)
		for _, hostName := range p.HostNames {
			splitted := strings.Split(hostName, ".")
			result = strings.Replace(result, hostName, p.Prefix+"-"+splitted[0]+"."+p.Prefix+"-"+splitted[1], -1)
		}
		return result
	}

	content.Accept(&data.Apply{Path: data.NewPath("metadata", "name"), Function: prefixName})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "serviceName"), Function: prefixName})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "env", ".*", "value"), Function: prefixHostName})

	content.Accept(&data.Apply{Path: data.NewPath("data", ".*"), Function: prefixHostName})

	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "envFrom", ".*", "configMapRef", "name"), Function: prefixName})

	labelsGetter := data.Get{Path: data.NewPath("metadata", "labels")}
	content.Accept(&labelsGetter)

}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Prefix",
			Description: "Add same prefix to all the k8s names",
			Parameter: []ProcessorParameter{
				{
					Name:        "prefix",
					Description: "The prefix to use before the name of the resources.",
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Prefix{}, config)
		},
	})
}
