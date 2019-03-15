package processor

import (
	"encoding/hex"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
)

import "crypto/md5"

type ConfigHash struct {
	DefaultProcessor
	Trigger    Trigger
	nameToHash map[string]string
}

func (processor *ConfigHash) Before(ctx *RenderContext, resources []*data.Resource) {
	processor.nameToHash = make(map[string]string)

	for _, resource := range resources {
		if resource.Kind() == "ConfigMap" && processor.Trigger.active(resource) {
			str, err := resource.Content.ToString()
			if err != nil {
				panic(err)
			}
			hash := md5.Sum([]byte(str))
			processor.nameToHash[resource.Name()] = hex.EncodeToString(hash[:md5.Size])
		}
	}

}

func (p *ConfigHash) BeforeResource(resource *data.Resource) {

	content := resource.Content
	getAll := data.GetAll{
		Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "envFrom", ".*", "configMapRef", "name"),
	}
	content.Accept(&getAll)
	for _, match := range getAll.Result {
		configName := match.Value.(*data.KeyNode).Value.(string)
		if val, ok := p.nameToHash[configName]; ok {
			annotations := data.SmartGetAll{
				Path: data.NewPath("metadata", "annotations"),
			}
			resource.Content.Accept(&annotations)
			for _, annotation := range annotations.Result {
				annotation.Value.(*data.MapNode).PutValue("flekszible/config-hash", val)
			}
			return
		}
	}

}
func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "ConfigHash",
			Description: "Add labels to the k8s resources with the hash of the used configmaps",
			Parameter: []ProcessorParameter{
				TriggerParameter,
			},
			Doc: `
Add a kubernetes annotation with the hash of the used configmap. With 
this approach you can force to re-create the k8s resources in case of config change. 
In case of configmap change the annotation value will be different and the resource
will be recreated.

As of now it supports only one configmap per resource and only the top-level
resource will be annotated (in case of statefulset this is the statefulset 
not the pod).

Example ('transformations/config.yaml'):
'''yaml
- type: ConfigHash
'''
` + TriggerDoc,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&ConfigHash{}, config)
		},
	})
}
