package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"strings"
)

type Run struct {
	DefaultProcessor
	Args    string
	Trigger Trigger
}

func (run *Run) BeforeResource(resource *data.Resource, location *ResourceNode) error {
	if run.Trigger.active(resource) {

		args := strings.Split(run.Args, " ")
		smartGet := data.SmartGetAll{Path: data.NewPath("spec", "template", "spec", "containers", ".*")}
		resource.Content.Accept(&smartGet)
		for _, result := range smartGet.Result {
			container := result.Value.(*data.MapNode)
			newArgs := container.CreateList("args")
			for _, arg := range args {
				newArgs.AddValue(arg)
			}
		}
	}
	return nil
}
func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Run",
			Description: "Replace args wi",
			Doc:         "Space separated string will be used as array",
			Parameters: []ProcessorParameter{
				TriggerParameter,
				{
					Name:        "args",
					Description: "Space separated arguments to use as args",
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Run{}, config)
		},
	})
}
