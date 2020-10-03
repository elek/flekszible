package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
)

type Mount struct {
	DefaultProcessor
	Trigger  Trigger
	HostPath string `yaml:"hostPath"`
	Path     string
	Name     string
}

func (mount *Mount) ToString() string {
	return CreateToString("mount").
		Add("path", mount.Path).
		Add("hostPath", mount.HostPath).
		Add("name", mount.Name).
		Build()
}

func (mount *Mount) BeforeResource(resource *data.Resource, location *ResourceNode) error {
	if mount.Trigger.active(resource) {

		volume := data.NewMapNode(data.RootPath())
		volume.PutValue("name", mount.Name)
		hostPath := volume.CreateMap("hostPath")
		hostPath.PutValue("path", mount.HostPath)

		addVolume := Add{
			Path:  data.NewPath("spec", "template", "spec", "volumes"),
			Value: data.ConvertToYaml(&volume),
		}

		volumeMount := data.NewMapNode(data.RootPath())
		volumeMount.PutValue("name", mount.Name)
		volumeMount.PutValue("mountPath", mount.Path)

		addVolumeMount := Add{
			Path:  data.NewPath("spec", "template", "spec", "(initC|c)ontainers", ".*", "volumeMounts"),
			Value: data.ConvertToYaml(&volumeMount),
		}

		err := addVolume.BeforeResource(resource, location)
		if err != nil {
			return err
		}

		err = addVolumeMount.BeforeResource(resource, location)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Mount",
			Description: "Mount external directory to the container (hostPath)",
			Doc:         "",
			Parameters: []ProcessorParameter{
				TriggerParameter,
				{
					Name:        "hostPath",
					Description: "Path on the host",
					Required:    true,
				},
				{
					Name:        "path",
					Description: "Path in the container",
					Required:    true,
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Mount{Name: "mount"}, config)
		},
	})
}
