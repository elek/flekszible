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

func (processor *ConfigHash) Before(ctx *RenderContext, resources []data.Resource) {
	processor.nameToHash = make(map[string]string)
	for _, resource := range resources {
		if resource.Kind() == "ConfigMap" && processor.Trigger.active(&resource) {
			str := ToString(&resource)
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
			Name: "ConfigHash",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&ConfigHash{}, config)
		},
	})
}
