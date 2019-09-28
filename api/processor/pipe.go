package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type Pipe struct {
	DefaultProcessor
	Command string
	Trigger Trigger
	Args    []string
}

func (p *Pipe) BeforeResource(resource *data.Resource) error {
	if !p.Trigger.active(resource) {
		return nil
	}
	converted := data.ConvertToYaml(resource.Content)
	bytes, err := yaml.Marshal(converted)
	if err != nil {
		return errors.Wrap(err, "Can't parse yaml file of "+resource.Path)
	}
	str := string(bytes)

	cmd := exec.Command(p.Command, p.Args...)

	builder := strings.Builder{}
	cmd.Stdin = strings.NewReader(str)
	cmd.Stdout = &builder
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Can't execute the command "+p.Command)
	}
	output := builder.String()
	logrus.Infof("Command is executed with the following output %s", output)
	node, err := data.ReadManifestString([]byte(output))
	if err != nil {
		return errors.Wrap(err, "Can't parse the result of the piped command for "+resource.Path)
	}
	resource.Content = node
	return nil
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Pipe",
			Description: "Transform content with external shell command.",
			Parameter: []ProcessorParameter{
				{
					Name:        "command",
					Description: "External program which transforms standard input to output",
				},
				{
					Name:        "args",
					Description: "List of the arguments of the command",
				},
				TriggerParameter,
			},
			Doc: `
Pipe executes a specific command to transform a k8s resources to a new resources.

The original manifest will be sent to the stdin of the process and the stdout will be processed as a the stdout of the file.

Example:

'''
- type: Pipe
  command: sed
  args: ["s/nginx/qwe/g"]
'''
` + TriggerDoc,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Pipe{}, config)
		},
	})
}
