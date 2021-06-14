package processor

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/elek/flekszible/api/v2/yaml"
)

type Kustomize struct {
	DefaultProcessor
}

func (writer *Kustomize) ToString() string {
	return "kustomize"
}

func (writer *Kustomize) Before(ctx *RenderContext, node *ResourceNode) error {
	_ = os.MkdirAll(ctx.OutputDir, 0755)
	destFile := path.Join(ctx.OutputDir, "kustomization.yaml")
	output := make(map[string]interface{})

	resources := make([]string, 0)

	for _, resource := range ctx.RootResource.AllResources() {
		if resource.Destination == "" {
			if resource.DestinationFileName != "" {
				resources = append(resources, resource.DestinationFileName)
			} else {
				resources = append(resources, CreateOutputFileName(resource.Name(), resource.Kind()))
			}
		}
	}
	output["resources"] = resources

	data, err := yaml.Marshal(output)
	if err != nil {
		return err
	}

	licenceHeader := ""
	licenceHeaderFile := path.Join(ctx.OutputDir, "LICENSE.header")
	if _, err := os.Stat(licenceHeader); os.IsNotExist(err) {
		content, _ := ioutil.ReadFile(licenceHeaderFile)
		licenceHeader = string(content) + "\n"
	}

	err = ioutil.WriteFile(destFile, append([]byte(licenceHeader), data...), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ActivateKustomize(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Kustomize",
			Description: "Generate Kustomize desciptor in the destination directory",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Kustomize{}, config)
		},
	})
}
